package handler

import (
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
	fmt.Println(chi.URLParam(r, "productID"))
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	checkoutParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(""),
		CancelURL:  stripe.String(""),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("__"),
				Quantity: stripe.Int64(1),
			},
		},
	}
	s, err := session.New(checkoutParams)
	if err != nil {
		return err
	}
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
	return nil
}
