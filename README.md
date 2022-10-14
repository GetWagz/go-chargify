# Chargify SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/GetWagz/go-chargify)](https://goreportcard.com/report/github.com/GetWagz/go-chargify)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)

This library is a small wrapper around a subset of the Chargify API.

This library is actively used in production, however not all end points are used regularly. We do our best to keep up to date with changes, but focus primarily on our own
needs. However, pull requests are always welcome!

## Usage

Usage is fairly straight-forward. See Configuration for more information about setting up and configuring the SDK.

## Environment Variables

* `CHARGIFY_ENV` set to production to actually make calls
* `CHARGIFY_API_KEY` Your secret API key
* `CHARGIFY_SUBDOMAIN` The subdomain for your account at Chargify

## Testing

Testing requires an actual account (nothing is mocked, but that could be a good addition!). Make sure your subdomain, api key, etc are properly set.

*IMPORTANT* If you run all of the tests, there isn't currently a way to delete the following entities. As such, you will need to handle that in the GUI until
a solution is provided in the official REST API:

* Product Family

## Moving CLI

A recent PR included a CLI. We will be moving it out into a separate repository. Do not rely on using it in this repo.

## Future Improvements

The following are improvements we would like to make, or would like PRs to address (reach out to us first to make sure no one else has started working on it though!):

* [ ] Migrate all `makeCall` and `makeEventsCall` invocations to the new `makeAPICall` function. This includes modifying the INTERNAL allocation to the options struct.

* [ ] Update calls to take options structs with pointers; add new functions where deprecation would occur so we can maintain backwards compatibility.

## Implemented Endpoints

Note that Chargify changes their API doc URLs regularly, so we have stopped providing links directly to the end points.

### Customers

* Create Customer
* Delete Customer
* Get Customers
* Search for Customers

### Events

* List Events
* List Events for Subscription
* Total Event Count
* Event Ingestion
* Bulk Event Ingestion

### Payment Profiles

* Create Payment Profile
* Delete Payment Profile

### Product Families

* Create Product Family
* Get Product Family
* Read Component By ID
* Read Component By Handle
* List Comonents for Product Family
* List Product Familiy via Site

### Products

* Create a Product
* Archive a Product
* Update a Product
* Get a Product By ID
* Get a Product By Handle
* Get a Product In Family

### Subscriptions

* Create Subscription
* Update Subscription
* Cancel Subscription - Immediately
* Cancel Subscription - Delayed
* Remove Delayed Cancellation
* List Subscriptions
* Purge a Subscription (only works in test mode!)

### Coupons

* Create a Coupon
* Find a Coupon
* Archive a Coupon
* List Coupons

## Hiring

Want to join an awesome team building cool products to improve the lives of pets and their owners? Send an email to engineering@wagz.com and let's find out if we're a good match! We are a remote-first company based in New Hampshire.

## Contributing

Pull Requests are welcome! See our `CONTRIBUTING.md` file for more information.
