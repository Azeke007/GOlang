func main() {
    p := Product{
        Name:     "Apple iPhone",
        Price:    999.99,
        Quantity: 10,
    }

    jsonBytes, err := productToJSON(p)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(string(jsonBytes))

    jsonStr := `{"Name":"Samsung TV","Price":1299.99,"Quantity":5}`
    product, err := jsonToProduct(jsonStr)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(product)
}