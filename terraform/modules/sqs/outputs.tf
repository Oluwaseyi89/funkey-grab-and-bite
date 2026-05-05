output "order_queue_url" {
  description = "Order queue URL"
  value       = aws_sqs_queue.order.url
}

output "order_queue_arn" {
  description = "Order queue ARN"
  value       = aws_sqs_queue.order.arn
}

output "order_queue_name" {
  description = "Order queue name"
  value       = aws_sqs_queue.order.name
}

output "catering_queue_url" {
  description = "Catering queue URL"
  value       = aws_sqs_queue.catering.url
}

output "catering_queue_arn" {
  description = "Catering queue ARN (Lambda event source)"
  value       = aws_sqs_queue.catering.arn
}

output "catering_queue_name" {
  description = "Catering queue name"
  value       = aws_sqs_queue.catering.name
}

output "order_dlq_arn" {
  description = "Order DLQ ARN"
  value       = aws_sqs_queue.order_dlq.arn
}

output "catering_dlq_arn" {
  description = "Catering DLQ ARN"
  value       = aws_sqs_queue.catering_dlq.arn
}
