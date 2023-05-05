# Set the API endpoint URL
$url = "http://localhost:8080/advertisements"

# Set the API request headers
$headers = @{
    "Content-Type" = "application/json"
}

# Set the API request body
$body = @(
    @{
        "name" = "Product A"
        "price" = 19.99
        "description" = "This is the description for Product A"
        "image" = "https://example.com/product-a.jpg"
        "author" = "John Doe"
        "date" = "2023-04-01T00:00:00Z"
    },
    @{
        "name" = "Product B"
        "price" = 29.99
        "description" = "This is the description for Product B"
        "image" = "https://example.com/product-b.jpg"
        "author" = "Jane Doe"
        "date" = "2023-04-05T00:00:00Z"
    },
    @{
        "name" = "Product C"
        "price" = 39.99
        "description" = "This is the description for Product C"
        "image" = "https://example.com/product-c.jpg"
        "author" = "Bob Smith"
        "date" = "2023-04-10T00:00:00Z"
    }
) | ConvertTo-Json

# Send 300 API requests
1..300 | ForEach-Object {
    Start-Job -ScriptBlock {
        Invoke-RestMethod -Method POST -Uri $using:url -Headers $using:headers -Body $using:body
    }
}

# Wait for all jobs to complete
Get-Job | Wait-Job | Receive-Job