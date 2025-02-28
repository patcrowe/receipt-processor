# receipt-processor
Receipt processor api challenge for Fetch Rewards

This is an api with two endpoints: '/receipts/process' and '/receipts/{id}/points'

/receipts/process takes a JSON receipt, calculates the points awarded for that receipt, stores those points in an in-memory map, then sends back the generated UUID for that receipt.
This will send back an error response if the receipt has invalid formatting or if elements of the receipt are invalid (ex: if the price of each item does not add up to the total).

/receipts/{id}/points retrieves the awarded points for the receipt with the UUID passed in through the {id} part of the url
This will send back an error response if the {id} passed in cannot be found in the receiptID-to-points map.

To run this program, navigate to the main receipt-processor/ directory and run the command ```go run .```

To send a receipt to the api, run the command: 
```curl -X POST 'http://localhost:8000/receipts/process' -d @receipFile.json```
In the above command, "receiptFile.json" represents the filename for the json file that contains your receipt data. An example in this repository would be Tests/test1.json
You can also write the json object directly in the command instead of using a file:
```curl -X POST 'http://localhost:8000/receipts/process' -d '{"retailer": "Target", "purchaseDate": "2022-01-02", "purchaseTime": "13:13", "total": "1.25", "items": [{"shortDescription": "Pepsi - 12-oz", "price": "1.25"}]}'```

To retrieve the points awarded for a receipt, copy the id that is returned from the process endpoint and use it as so:
```curl -X GET 'http//localhost:8000/receipts/Id-Goes-Here/points'```

Using Tests/test1.json as an example:
```curl -X POST 'http://localhost:8000/receipts/process' -d @Tests/test1.json```
Returns {"id":"32ff3d40-3c12-57e4-ab1f-30a33864be7f"}
```curl -X GET 'http://localhost:8000/receipts/32ff3d40-3c12-57e4-ab1f-30a33864be7f/points'```
Returns {"points":31}