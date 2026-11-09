package repository

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Runs only with TEST_DATABASE_URL; see openIntegrationDB.
func TestAccountCodesWorkOnceAndLockAfterWrongTries(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()
	tag := fmt.Sprint(time.Now().UnixNano())
	var userID int
	if err := db.QueryRow(`INSERT INTO users (phone, full_name, password_hash) VALUES ($1, 'Code test', '') RETURNING id`, "0809"+tag[len(tag)-7:]).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM users WHERE id = $1`, userID) })
	codes := NewAccountCodeRepository(db)

	use := func(hash string) bool {
		t.Helper()
		ok, err := codes.Use(ctx, userID, CodePurposeClaimAccount, hash, 5)
		if err != nil {
			t.Fatal(err)
		}
		return ok
	}

	_ = codes.Create(ctx, userID, CodePurposeClaimAccount, "a@b.c", "first", time.Minute)
	_ = codes.Create(ctx, userID, CodePurposeClaimAccount, "a@b.c", "second", time.Minute)
	if use("first") {
		t.Fatal("an earlier code must stop working once a new one is sent")
	}
	if !use("second") || use("second") {
		t.Fatal("the latest code must work exactly once")
	}

	_ = codes.Create(ctx, userID, CodePurposeClaimAccount, "a@b.c", "third", time.Minute)
	for range 5 {
		use("wrong")
	}
	if use("third") {
		t.Fatal("after five wrong tries even the right code must be refused")
	}

	_ = codes.Create(ctx, userID, CodePurposePasswordReset, "a@b.c", "expired", -time.Second)
	if ok, _ := codes.Use(ctx, userID, CodePurposePasswordReset, "expired", 5); ok {
		t.Fatal("an expired code must be refused")
	}
}

func TestUsersMatchAnyPhoneFormatAndClaimOnce(t *testing.T) {
	db := openIntegrationDB(t)
	tag := fmt.Sprint(time.Now().UnixNano())
	local := "0808" + tag[len(tag)-7:]
	users := NewUserRepository(db)
	var guestID int
	if err := db.QueryRow(`INSERT INTO users (phone, full_name, password_hash, email) VALUES ($1, 'Format test', '', $2) RETURNING id`,
		"+234"+local[1:], "fmt"+tag+"@example.test").Scan(&guestID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec(`DELETE FROM users WHERE id = $1`, guestID) })

	for _, typed := range []string{local, "+234 " + local[1:4] + " " + local[4:], "234" + local[1:]} {
		found, err := users.FindByPhoneAnyFormat(typed)
		if err != nil || found == nil || found.ID != guestID {
			t.Fatalf("%q should find the customer, got %+v %v", typed, found, err)
		}
	}
	if found, _ := users.FindByEmail("FMT" + tag + "@EXAMPLE.TEST"); found == nil || found.ID != guestID {
		t.Fatal("email lookup must ignore case")
	}

	if ok, err := users.ClaimGuest(guestID, "Format Owner", "fmt"+tag+"@example.test", "hash"); err != nil || !ok {
		t.Fatalf("expected the guest record claimed, got %v %v", ok, err)
	}
	if ok, _ := users.ClaimGuest(guestID, "Someone else", "x@example.test", "other"); ok {
		t.Fatal("a claimed record can't be claimed again")
	}
}

func TestGuestOrdersOnlyLinkToAccountsWithMatchingEmail(t *testing.T) {
	f := newCheckoutFixture(t, 10, 10)
	ctx := context.Background()
	local := "0807" + f.tag[len(f.tag)-7:]
	accountEmail := "owner" + f.tag + "@example.test"
	var accountID int
	if err := f.db.QueryRow(`INSERT INTO users (phone, full_name, password_hash, email) VALUES ($1, $2, 'hash', $3) RETURNING id`,
		"+234"+local[1:], "Checkout guest "+f.tag+" owner", accountEmail).Scan(&accountID); err != nil {
		t.Fatal(err)
	}

	stranger := "someone" + f.tag + "@example.test"
	unlinked, err := f.repo.PlaceOrder(ctx, f.checkout(1, 1, local, &stranger, false))
	if err != nil {
		t.Fatal(err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM orders WHERE id = $1 AND user_id IS NULL`, unlinked); n != 1 {
		t.Fatal("a guest order with another email must not be linked to the account owning the number")
	}

	linked, err := f.repo.PlaceOrder(ctx, f.checkout(2, 1, local, &accountEmail, false))
	if err != nil {
		t.Fatal(err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM orders WHERE id = $1 AND user_id = $2`, linked, accountID); n != 1 {
		t.Fatal("a guest order with the account's own email links to it, whatever the number's format")
	}

	// A guest record is matched across formats instead of being duplicated.
	guestLocal := "0806" + f.tag[len(f.tag)-7:]
	if _, err := f.repo.PlaceOrder(ctx, f.checkout(3, 1, guestLocal, nil, false)); err != nil {
		t.Fatal(err)
	}
	if _, err := f.repo.PlaceOrder(ctx, f.checkout(4, 1, "+234"+guestLocal[1:], nil, false)); err != nil {
		t.Fatal(err)
	}
	if n := f.count(t, `SELECT COUNT(*) FROM users WHERE right(regexp_replace(phone, '\D', '', 'g'), 10) = $1`, guestLocal[1:]); n != 1 {
		t.Fatalf("expected one customer record for one number in two formats, got %d", n)
	}
}
