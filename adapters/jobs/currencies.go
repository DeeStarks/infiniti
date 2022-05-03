package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deestarks/infiniti/application/services"
	"github.com/deestarks/infiniti/config"
	"github.com/deestarks/infiniti/utils"
	"github.com/go-co-op/gocron"
)

type ExchangeRateUpdate struct {
	Add			bool
	appPort		services.AppServicesPort
	scheduler	*gocron.Scheduler
}

func (adapter *JobAdapter) ExchangeRateUpdate() ExchangeRateUpdate {
	// Use the API key to determine if the job should be added to the job queue
	var hasKey = true
	if config.GetExchangeRateAPIKey() == "" {
		hasKey = false
		utils.LogJobMessage("Error: Exchange rate API key not set")
	}

	return ExchangeRateUpdate{
		Add: hasKey,
		appPort: adapter.appPort,
		scheduler: adapter.scheduler,
	}
}

type ExchangeRateResource struct {
	Status 	string 				`json:"result"`
	Rates	map[string]float64 	`json:"conversion_rates"`
}

func (job ExchangeRateUpdate) Run() {
	job.scheduler.Cron("0 6 * * *").Do(func() { // Run every day at 6:00 AM
		var rates ExchangeRateResource

		// Get the exchange rates
		res, err := http.Get(fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", config.GetExchangeRateAPIKey()))
		if err != nil {
			utils.LogJobMessage("Error: %s", err.Error())
			return
		}

		err = json.NewDecoder(res.Body).Decode(&rates)
		if err != nil {
			utils.LogJobMessage("Error: %s", err.Error())
			return
		}

		// Fetch all currencies
		service := job.appPort.NewCurrencyService()
		currencies, err := service.ListCurrencies()
		if err, ok := err.(*utils.RequestError); ok {
			utils.LogJobMessage("Error: %s", err.Error())
			return
		}

		for _, currency := range currencies {
			// Update the currency
			if _, ok := rates.Rates[currency.Symbol]; ok {
				_, err = service.UpdateCurrency("id", currency.Id, map[string]interface{}{
					"conversion_rate_to_usd": rates.Rates[currency.Symbol],
				})
				if err, ok := err.(*utils.RequestError); ok {
					utils.LogJobMessage("Error: %s", err.Error())
					return
				}
			}
		}
		utils.LogJobMessage("Exchange rates updated")

	})
	job.scheduler.StartAsync()
}