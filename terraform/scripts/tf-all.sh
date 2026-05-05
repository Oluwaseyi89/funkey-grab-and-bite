#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ENVS=(staging production)
TMP_PLAN_ROOT=""

cleanup() {
  if [[ -n "${TMP_PLAN_ROOT:-}" && -d "${TMP_PLAN_ROOT:-}" ]]; then
    rm -rf "$TMP_PLAN_ROOT"
  fi
}

set_tf_var_defaults() {
  if [[ -z "${TF_VAR_alert_email:-}" ]]; then
    if [[ -n "${alert_email:-}" ]]; then
      export TF_VAR_alert_email="${alert_email}"
    elif [[ -n "${ALERT_EMAIL:-}" ]]; then
      export TF_VAR_alert_email="${ALERT_EMAIL}"
    fi
  fi

  if [[ -z "${TF_VAR_ses_sender_email:-}" ]]; then
    if [[ -n "${ses_sender_email:-}" ]]; then
      export TF_VAR_ses_sender_email="${ses_sender_email}"
    elif [[ -n "${SES_SENDER_EMAIL:-}" ]]; then
      export TF_VAR_ses_sender_email="${SES_SENDER_EMAIL}"
    fi
  fi

  if [[ -z "${TF_VAR_api_image_tag:-}" ]]; then
    if [[ -n "${api_image_tag:-}" ]]; then
      export TF_VAR_api_image_tag="${api_image_tag}"
    elif [[ -n "${API_IMAGE_TAG:-}" ]]; then
      export TF_VAR_api_image_tag="${API_IMAGE_TAG}"
    fi
  fi
}

tf() {
  local configured_region
  local -a env_args

  configured_region="$(aws configure get region 2>/dev/null || true)"
  env_args=("AWS_SDK_LOAD_CONFIG=1")

  if [[ -n "${AWS_PROFILE:-}" ]]; then
    env_args+=("AWS_PROFILE=${AWS_PROFILE}")
  fi

  if [[ -n "${AWS_REGION:-}" ]]; then
    env_args+=("AWS_REGION=${AWS_REGION}")
  elif [[ -n "${AWS_DEFAULT_REGION:-}" ]]; then
    env_args+=("AWS_REGION=${AWS_DEFAULT_REGION}")
  elif [[ -n "$configured_region" ]]; then
    env_args+=("AWS_REGION=${configured_region}")
  fi

  env "${env_args[@]}" terraform "$@"
}

usage() {
  cat <<'EOF'
Usage:
  ./scripts/tf-all.sh fmt
  ./scripts/tf-all.sh validate
  ./scripts/tf-all.sh plan

Behavior:
  fmt      Runs terraform fmt recursively from terraform/.
  validate Runs init -backend=false -reconfigure and validate in staging/production.
  plan     Runs init -backend=false -reconfigure and plan -lock=false in staging/production.
EOF
}

run_fmt() {
  echo "==> terraform fmt -recursive"
  tf -chdir="$ROOT_DIR" fmt -recursive
}

run_validate() {
  for env_name in "${ENVS[@]}"; do
    echo "==> [$env_name] terraform init -backend=false -reconfigure"
    tf -chdir="$ROOT_DIR/environments/$env_name" init -backend=false -reconfigure -input=false

    echo "==> [$env_name] terraform validate"
    tf -chdir="$ROOT_DIR/environments/$env_name" validate
  done
}

run_plan() {
  local rc
  local -a plan_args

  TMP_PLAN_ROOT="$(mktemp -d /tmp/funkey-tf-plan-XXXXXX)"
  trap cleanup EXIT

  cp -r "$ROOT_DIR"/* "$TMP_PLAN_ROOT"/
  rm -f "$TMP_PLAN_ROOT"/environments/staging/backend.tf
  rm -f "$TMP_PLAN_ROOT"/environments/production/backend.tf
  rm -rf "$TMP_PLAN_ROOT"/environments/staging/.terraform
  rm -rf "$TMP_PLAN_ROOT"/environments/production/.terraform

  echo "INFO: using temporary backend-free workspace at: $TMP_PLAN_ROOT"
  set_tf_var_defaults

  for env_name in "${ENVS[@]}"; do
    plan_args=(-input=false -lock=false -refresh=true -detailed-exitcode)
    if [[ -f "$TMP_PLAN_ROOT/environments/$env_name/terraform.tfvars" ]]; then
      plan_args+=("-var-file=terraform.tfvars")
    elif [[ -f "$TMP_PLAN_ROOT/environments/$env_name/terraform.tfvars.example" ]]; then
      plan_args+=("-var-file=terraform.tfvars.example")
    fi

    echo "==> [$env_name] terraform init -backend=false -reconfigure"
    tf -chdir="$TMP_PLAN_ROOT/environments/$env_name" init -backend=false -reconfigure -input=false

    echo "==> [$env_name] terraform plan"
    set +e
    tf -chdir="$TMP_PLAN_ROOT/environments/$env_name" plan "${plan_args[@]}"
    rc=$?
    set -e

    if [[ $rc -eq 1 ]]; then
      echo "ERROR: plan failed for environment: $env_name"
      exit 1
    elif [[ $rc -eq 2 ]]; then
      echo "INFO: plan succeeded with pending changes for environment: $env_name"
    else
      echo "INFO: plan succeeded with no changes for environment: $env_name"
    fi
  done
}

if [[ $# -ne 1 ]]; then
  usage
  exit 1
fi

case "$1" in
  fmt)
    run_fmt
    ;;
  validate)
    run_validate
    ;;
  plan)
    run_plan
    ;;
  *)
    usage
    exit 1
    ;;
esac
