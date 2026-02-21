# Paramount+

## paypal.com

1. about:config
2. general.useragent.override
3. string
4. add
5. Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:147.0) Gecko/20100101 Firefox/147.0
6. paramountplus.com
7. get started
8. paramount+ premium
   - continue
9. full name
10. email
   - mail.tm
11. password
12. zip code
13. birthdate
14. gender
15. agree & continue
16. paypal
17. continue to paypal
18. agree and continue
19. subscribe
20. paypal.com/myaccount/autopay
21. paramount
22. stop paying with paypal

## privacy.com

1. about:config
2. general.useragent.override
3. string
4. add
5. Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:147.0) Gecko/20100101 Firefox/147.0
6. paramountplus.com
7. get started
8. paramount+ essential
   - continue
9. full name
10. email
11. password
12. zip code
13. birthdate
14. gender
15. agree & continue
16. first name
17. last name
18. address
19. city
20. state
21. zip
22. credit card
23. exp MM
24. YYYY
25. CVV
26. subscribe

~~~
POST /account/xhr/processPayment/ HTTP/2
Host: www.paramountplus.com

HTTP/2 200 OK

message = "Error occurred. Please try again.";
recaptchaTokenValidated = true;
success = false;
trackingMessage = "Could not purchase subscription. Error:
ErrorResponse(error=ErrorResponse.ErrorDetail(type=transaction, message=The
transaction was declined. Please use a different card, contact your bank, or
contact support., params=null,
transactionError=ErrorResponse.TransactionError(object=transaction_error,
transactionId=ybs6siqfsn9x, category=fraud, threeDSecureActionTokenId=null,
code=fraud_gateway)))";
~~~
