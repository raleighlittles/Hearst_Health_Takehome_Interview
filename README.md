# Assignment

Hearst is building an online marketplace where retailers can offer products for sale that are of interest to readers of our newsletter. Many products will be offered on this marketplace, and many products will have offerings from multiple retailers.
We want to provide our subscribers with the current lowest price for the products they are interested in. 
To achieve this, we are building
a service that can quickly respond with the lowest price for items by [SKU](https://en.wikipedia.org/wiki/Stock_keeping_unit) for a product detail page or an in-place advertisement.
The system provides two main entry points. 

The webhook `receive()` is called any time a new price is made available from one of
our known retailers. These can occur at any time. The receiver should ensure we store the price received for 2 purposes:
1. Vending the current lowest price on a product across retailers
2. performing analysis on the price history of an item.

The method signature is as follows:

```typescript
interface PriceUpdate {
retailer: string; // the name of the retailer
sku: string; // assume retailers share a common SKU
price: number; // always the price per unit
url?: string; // product detail link, optional
}
```

```typescript
function receive(payload: PriceUpdate): void;
```

The `findPrice()` endpoint is responsible for determining the best price of an item at the time it is called. It is exposed to web-
scale traffic, so must be able to respond in under 20ms. The method signature and response payload are as follows:

```typescript
interface ProductPrice {
retailer: string;
sku: string;
price: number;
url?: string;
}
```

```typescript
function findPrice(sku: string): ProductPrice;
```

Your job is to implement the webhook `receive()` and `findPrice()` functions in a language of your choosing. 
Design and describe your
persistence strategy including schemas, indexes, and rationale. 
You do not need to implement the persistence methods, only add a
comment where it should be used. Please document any third-party libraries and any assumptions you have made in your
comments.

# Sample data

| Id | SKU   | Retailer | Price |
|----|-------|----------|-------|
| 1  | CLOCK | Walmart  | $20   |
| 2  | BED   | IKEA     | $140  |
| 3  | CLOCK | Target   | $15   |
| 4  | CLOCK | Target   | $14   |
| 5  | CLOCK | Best Buy | $30   |
| 6  | BED   | Wayfair  | $120  |
| 7  | CLOCK | Target   | $25   |
| 8  | CLOCK | Walmart  | $27   |
| 9  | CLOCK | Costco   | $12   |
| 10 | BED   | IKEA     | $100  |
| 11 | CLOCK | Costco   | $13   |

# Sample inputs and outputs

| After insertion of this row | findPrice() Input | findPrice() Output - [Retailer, Low Price] |   |   |
|-----------------------------|-------------------|--------------------------------------------|---|---|
| 1                           | CLOCK             | [Walmart, $20]                             |   |   |
| 2                           | BED               | [IKEA, $140]                               |   |   |
| 3                           | CLOCK             | [Target, $15]                              |   |   |
| 4                           | CLOCK             | [Target $14]                               |   |   |
| 5                           | CLOCK             | [Target, $14]                              |   |   |
| 6                           | BED               | [WayFair, $120]                            |   |   |
| 7                           | CLOCK             | [Walmart, $20]                             |   |   |
| 8                           | CLOCK             | [Target, $25]                              |   |   |
| 9                           | CLOCK             | [Costco, $12]                              |   |   |
| 10                          | BED               | [IKEA, $100]                               |   |   |
| 11                          | CLOCK             | [Costco, $13]                              |   |   |
