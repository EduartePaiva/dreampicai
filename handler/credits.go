package handler

import (
	"dreampicai/db"
	"dreampicai/view/credits"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, credits.Index())
}

func HandleStripeCheckoutCreate(w http.ResponseWriter, r *http.Request) error {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:3000/checkout/success/{CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:3000/checkout/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(chi.URLParam(r, "productID")),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}
	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}
	return hxRedirect(w, r, s.URL)
}

func HandleStripeCheckoutSuccess(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	sessionID := chi.URLParam(r, "sessionID")
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		return err
	}
	lineItemParams := stripe.CheckoutSessionListLineItemsParams{}
	lineItemParams.Session = stripe.String(sess.ID)
	iter := session.ListLineItems(&lineItemParams)
	iter.Next()
	item := iter.LineItem()
	priceID := item.Price.ID
	switch priceID {
	case os.Getenv("STRIPE_100_CREDITS"):
		user.Account.Credits += 100
	case os.Getenv("STRIPE_250_CREDITS"):
		user.Account.Credits += 250
	case os.Getenv("STRIPE_500_CREDITS"):
		user.Account.Credits += 500
	default:
		return fmt.Errorf("invalid price id: %s", priceID)
	}
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	http.Redirect(w, r, "/generate", http.StatusSeeOther)
	return nil
}
func HandleStripeCheckoutCancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}
