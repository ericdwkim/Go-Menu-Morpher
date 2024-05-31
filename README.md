# 🌟 Go, Menu Morpher! 🌟

![img.png](asset/go_menu_morpher.png)
## Description

Go, Menu Morpher! is the famous [Menu Morpher](https://github.com/ericdwkim/Menu-Morpher) Python client rewritten in Golang for the Google My Business API. 

This tool allows businesses to easily manage and transform their menu data programmatically. No more hassle in finding your `account_id`, `location_id`, and reading Google API documentation! 

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go version 1.22.3 or higher
- Completed the [Prerequisites - Google Business Profile APIs](https://developers.google.com/my-business/content/prereqs)
- Have your `project_id` handy
- [OAuth2.0 credentials](https://developers.google.com/my-business/content/basic-setup#request-client-id) requested, created, and [consented](https://console.cloud.google.com/apis/credentials/consent?project={your_project_id_here})
- Enabled the following (3) APIs via `APIs & Services` through Google Cloud Console (GCC)

NOTE: To use the following hyperlinks, replace `{your_project_id_here}` with your _actual_ `project_id`:

### [GCC - APIs & Services](https://console.cloud.google.com/apis/dashboard?project={your_project_id_here})

#### [MyBusinessBusinessInformation](https://console.cloud.google.com/apis/api/mybusinessbusinessinformation.googleapis.com/metrics?project={your_project_id_here})
serviceName:`mybusinessbusinessinformation` |
version: `v1`
#### [MyBusinessAccountManagement](https://console.cloud.google.com/apis/api/mybusinessaccountmanagement.googleapis.com/metrics?project={your_project_id_here})
serviceName:`mybusinessaccountmanagement` |
version: `v1`
#### [Google My Business](https://console.cloud.google.com/apis/api/mybusiness.googleapis.com/quotas?project={your_project_id_here}) 
serviceName: `mybusiness` (aka "Google My Business") |
version: `v4`


## Installation

### Clone the Repository

Start by cloning the repository to your local machine:

```bash
git clone https://github.com/ericdwkim/Go-Menu-Morpher.git
cd Go-Menu-Morpher
```
