# disney

here is what the web client does, note we can probably omit some of these calls.
first it does `registerDevice`:

~~~
POST https://disney.api.edge.bamgrid.com/graph/v1/device/graphql HTTP/2.0
authorization: Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu...

{
  "query": "mutation registerDevice($input: RegisterDeviceInput!) {\n      registerDevice(registerDevice: $input) {\n        grant {\n          grantType\n          assertion\n        },\n        token {\n          accessToken\n          accessTokenType\n          expiresIn\n          refreshToken\n          tokenType\n        },\n        session: activeSession {\n          sessionId\n          partnerName\n          device {\n            id\n            category\n            platform\n          }\n          profile {\n            id\n          }\n          experiments {\n            featureId\n            variantId\n            version\n          }\n          portabilityLocation {\n            countryCode\n            type\n          }\n          homeLocation {\n            adsSupported\n            countryCode\n          }\n          household {\n            householdScore\n          }\n          preferredMaturityRating {\n            impliedMaturityRating\n            ratingSystem\n          }\n          identity {\n            id\n          }\n          location {\n            adsSupported\n            type\n            countryCode\n            dma\n            asn\n            regionName\n            connectionType\n            zipCode\n          }\n        }\n      }\n    }",
  "variables": {
    "input": {
      "deviceProfile": "windows",
      "deviceFamily": "browser",
      "applicationRuntime": "firefox",
      "attributes": {
        "operatingSystem": "windows",
        "operatingSystemVersion": "10.0"
      }
    }
  }
}

HTTP/2.0 200 

{
  "data": {
    "registerDevice": {
      "token": {
        "accessTokenType": "Device",
        "refreshToken": "...BzU_HikpzuPbDyTvaXpzDRmxS0n1NqR7e20tEjoJSfirpos-..."
      }
    }
  }
}
~~~

use `refreshToken` for next request:

~~~
POST https://disney.api.edge.bamgrid.com/graph/v1/device/graphql HTTP/2.0
authorization: Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu...

{
  "query": "mutation refreshToken($input:RefreshTokenInput!){refreshToken(refreshToken:$input){activeSession{sessionId}}}",
  "variables": {
    "input": {
      "refreshToken": "...BzU_HikpzuPbDyTvaXpzDRmxS0n1NqR7e20tEjoJSfirpos-RFr..."
    }
  },
  "operationName": "refreshToken"
}

HTTP/2.0 200 

{
  "extensions": {
    "sdk": {
      "token": {
        "accessTokenType": "Device",
        "accessToken": "eyJ6aXAiOiJERUYiLCJraWQiOiJ0Vy10M2ZQUTJEN2Q0YlBWTU1rSkd4dkJlZ0ZXQkdXek5KcFFtOGRJMWYwIiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..j8By7jF572d4uGmX.8HGabdVXIQpO0WiHeyJguod1Mwf0SgS1WMI2al-Ehz8ymzqOrFvYRtyQXccuSuIF8yST_XG6haodpHiG6bzHO3Ft8NXoiNAmYavMTXN4aDinARfoxtwpsbJROprsNt2Qx5Xn4T_5rq_NxjAfpqKFWBbAkXGOF33H5aVmq-GY3G9Rp4OdhWX7_VDZI5EzakRznuOt3CSaituPTq7d29ZzOJqOQNXIoVBMYXk16deMPLBHQbLlrjjewdAKX4EBz9c_qX4zONJ897YcTJnL-LDpX3l_XuvsgZRDWFF4EUrQoI_wFwtO4X1RJMD8IMxG_GM0S_Uzt1j9OKsAwFs-a9JxLHqiFxBj5NKacWESKSbDpN7UZviPdLvL3Pf3v1Kvh3KSjewabs6VY489O2bziKrznwVls_94yD9MCJ4-vEQAicGvqO7YERDDvzT3dtCnA6RJL6_jbt0Br3KGjIqHJcpfj43irpPvNfWtMtKf_auL0lUA9yfJm2ijoH69OZfTsgZIM8MZfCsCZN6x09bKXRI8k4AVn6lTGBzVCa9Qkb7qGBpy8FzK6ePvWY4x8RX_EBBnlh7yoj9tq8WD3V72oGKlvuh2TSt0H2z86GgG31GV9SAdiMvShVvpnEbCNQbcL8G5gV9KNqudMY12nuTn_5chkzP8IJq6y7gTb1qqSAujaGABwAm0pkAcbvL7-xZ7LV7gffRv5JmKbLk5yR3X3KasM0SaUrVw5N5pG8WTosCTYf7KRizMxTSbpwic_UIc5rTr6cl1ftE4Nvj4xqWxn7NEZ2tAuPeVLQCqebFfRLvb_5VTHQHUYVpt5ulFQ7rEwZoES_-rBz9bxDE1IY0mwIaFHX4RXRloDmSbAXzm3w3Oz7D0cQDyC7e17nWyGFbRZjQ0-1RmRs0FzyFOC4-6T4VXP_ZO-F5lbwMSFRACmlGY3ZjRtQeEhqLnwuIoCvD7LPjSVlp5YrPirLE3LBJEjtxT6WDEPadAn8zNlLB1Y0kH1nXdaI6PpsymOaxEu7xuo5gDYESwrvKcRm3ZhauLKEma7q3q1-uXZ5DiNPlD4qJ6Amj5X77qRxN_h7n8eOWuOkD2MrhWShDMGqdEE4CYIpzcHM3-GkiLIhkfmsyXBj_iv2S2U1vbqUPJH2MY-QubEJXoUJIpoMAbnKT95wUL1ml8P1mp-kJi1MJyc7LiQuJXdFvxtF0a2Sz_X4L0RAGiZa6MDeh_C5BE8v3WgnBZ3Sa9UEjnRBTh3__Ru1_WpRPGlbimj4VNcXK8ch3MEZSmB8CoW28zcm8ZyUCQcVSf0WYC7GSL3lCnHmUyxMWLBweE2Xn_iAG99hzuPPn1vYaEpibKa-WswOw4ObaYK_Spbt5-IgJ85ZxgkunAOZ96k1D_XagB52cILu74qJzn1Cs9k2-ARlz-8q989PUvAkkp3jF9tZam97U4ZJqVeyHEVAWKePeP6wY7trGkq0dXpi8bK12KpOa65j7Tj4Jk_n7F7OvB95SfA5Mvlyg2mXTqXjAJwrbmM2JZG1EgafUhjXVa1kSh5I1E7Khw0MOtvomcoY79lA5PyGnrX0L5rQLK3wuI78tjEu4vZl6Q5YINiURC5A8qsvYDZ15J2KnDONFn2D8qDl-AkEG0YuFcSaSeWpToID_Q7Z1KO8JTXEpKsIQ07Tsjfg_g7DRzSNOKxYQ75sFiVrT_tN0Rcz1Y3dP4Qhv8ur8npnf5CdrDB2b9V2RN_edV6KIv8qZw2E8NXCjnelzoEOrAm9UwnB2vt_GNETGil36cvx2r2MnJLyp8Z2of1i9zr_23HqmmvdHUCHZJ6SiMzEPILtOvpTBVcM50Y18aTu-yDQG1VfFuxNOQEPJfiymnCvxxHW2nWFs0UVBFjJi4RX6t-cA2B4YdD62wRy0DMu18d-tUW9qxiD1-OSzJEUKB1dnuji_53yP-6SvE9sI6S7nRb0vlzb1PNmgpjiLTcwfDuqwBrJYW4Iy4QYID98XDXMM-2WJtStIq3V9Grs107CAvnC5V6vaxWirkYiKCeEI9nKrxqy_E8oqOWf93S9rjHH_N7E3ufFN6BmoD2Xe8ZUHPCc9oJB8G0CCBHgwI98w_5qrP5J3eeBE8AuDN2c88GOIHuTcMRE-nQ8EdodvFCwV0NHwJZgKDlZFyqIqzmP_C2PRqRRNGoKNwLsKhVSPGpYg8obhAdESAYylmCtA_Tb7vQwMHHXgOcQhXnVrV3Dc9Gk39l9kdZlby68hDTVBd-H8jJVirjKQ5SvLxf3PdMoTb5lPYQLOrVd5KjoZYR6KIflLskp-0x2pUNrzk6cv8Oq-krJipASJYFIjzm66FJQ8A8MhnMwMuRx2aM--w9iPCJYaAEDS93QdnkKV43y1WASJylg6i0J5JCLhAxtWz-nwZgY3YxCeTA65A4n1lZIEJiXN4LuK9aR95ieQ93TGpISVQKVa3s7BKrJfSWjbn1s6cfu3ozJMBCeoDTkOtunkjdEhhn2A0XEg2IfMg7Btte1ElXtGZOvlZycb01nYuwjSd3MVz6kmVlpKIxpFjftpx0W5nerMcEd6lm1nPt_PwG6WTjLnthteK7fSEPwnKlsSd39HSh1JVzMj0Nkb7b-dRJjQffEOa1p8Ct94O3qL1gRSpIrlY0I06IQ3VjV0_m6J4__U.di-JNIDbdaEzpX35aKnuRQ"
      }
    }
  }
}
~~~

use `accessToken` with next request:

~~~
POST https://disney.api.edge.bamgrid.com/v1/public/graphql HTTP/2.0
authorization: Bearer ...8HGabdVXIQpO0WiHeyJguod1Mwf0SgS1WMI2al-Ehz8ymzqOrFvYR...

{
  "query": "\n    mutation login($input: LoginInput!) {\n        login(login: $input) {\n            actionGrant\n            account {\n              activeProfile {\n                id\n              }\n              profiles {\n                id\n                attributes {\n                  isDefault\n                  parentalControls {\n                    isPinProtected\n                  }\n                }\n              }\n            }\n            activeSession {\n              isSubscriber\n            }\n            identity {\n              personalInfo {\n                dateOfBirth\n                gender\n              }\n              flows {\n                personalInfo {\n                  requiresCollection\n                  eligibleForCollection\n                }\n              }\n            }\n        }\n    }\n",
  "variables": {
    "input": {
      "email": "EMAIL",
      "password": "PASSWORD"
    }
  },
  "operationName": "login"
}
~~~

current process:

1. `registerDevice`
2. `refreshToken`
3. login

would this work:

1. `registerDevice`
2. login
