package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiKey = "b316e369-3bb8-4177-80ca-7b692ed0eb1c"
const finnHubApiKey = "d16bumpr01qvtdbhkk5gd16bumpr01qvtdbhkk60"

func fetchCryptoPrice(symbol string) string {
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s", symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//required for coinmarketcap API calls
	req.Header.Add("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Optional: print raw JSON to see structure
	//fmt.Println(string(body))

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Navigate through the JSON to get the price
	// Example: result["data"].(map[string]interface{})["BTC"].(map[string]interface{})["quote"].(map[string]interface{})["USD"].(map[string]interface{})["price"]

	data := result["data"].(map[string]interface{})[symbol].(map[string]interface{})
	quote := data["quote"].(map[string]interface{})["USD"].(map[string]interface{})
	price := quote["price"].(float64)

	formattedString := fmt.Sprintf("%s: $%.2f\n", symbol, price)
	return formattedString
}

func fetchStockPrice(symbol string) string {
	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s", symbol)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Finnhub-Token", finnHubApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	error := json.Unmarshal(body, &result)
	if error != nil {
		log.Fatal("Error unmarshalling JSON:", error)
	}

	priceVal, ok := result["c"]
	if !ok || priceVal == nil {
		return fmt.Sprintf("No valid price found for %s", symbol)
	}

	return fmt.Sprintf("%s: $%.2f\n", symbol, priceVal)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi welcome to this portoflio tracker!")
	fmt.Fprintln(w, "Go to /crypto to see the current price of crypto in your portfolio!")
	fmt.Fprintln(w, "Go to /stock to see the current price of stock in your portfolio!")
}

func cryptoHandler(w http.ResponseWriter, r *http.Request) {
	bitcoin := fetchCryptoPrice("BTC")
	etherium := fetchCryptoPrice("ETH")
	fmt.Fprintln(w, bitcoin)
	fmt.Fprintln(w, etherium)
}

func stockHandler(w http.ResponseWriter, r *http.Request) {
	VanguardETF := fetchStockPrice("VOO")
	InvescoETF := fetchStockPrice("QQQ")
	fmt.Fprintln(w, VanguardETF)
	fmt.Fprintln(w, InvescoETF)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/crypto", cryptoHandler)
	http.HandleFunc("/stock", stockHandler)

	fmt.Println("Server started at port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Couldn't host at :8080")
	}
}
