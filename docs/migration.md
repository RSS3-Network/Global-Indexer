# API Migration Guide:

## Background

In 2022, we launched the [PreGod](https://generator.swagger.io/?url=https%3A%2F%2Fpregod.rss3.dev%2Fv1%2Fopenapi%3Fjson%3Dtrue#/) version of our open information indexing service. This initial version was built on a relatively complex caching architecture, requiring users to manually trigger cache refreshes to access the latest data. This process was not only cumbersome but also less efficient in handling large volumes of data requests.

To address these issues, we have developed the Node version, which abandons the previous caching architecture in favor of a more efficient and powerful indexing mechanism. The new architecture enhances the performance and scalability of our services significantly. The [Node API]( https://developer.dev.rss3.io/reference) simplifies the process of data retrieval, enabling users to receive updated data directly through API calls without the need for manual interventions.

The shift to this new system involves substantial changes to the APIs, necessitating a detailed migration plan to ensure a smooth transition for our users. Below, we outline the necessary steps and considerations for migrating from the PreGod API to the Node API.


## API Migration Explanation

### 1. GET `/notes/{address}` -> GET `/decentralized/{account}`

#### Endpoint Summary
This API transition is from `GET /notes/{address}` in the PreGod API to `GET /decentralized/{account}` in the Node API. Both endpoints serve a similar purpose: to retrieve open information activity records for a specified account. The transition to the Node API represents an evolution in data access and querying capabilities, with a more streamlined and efficient set of parameters tailored for modern use cases.

#### Brief Description of the API

Both endpoints are designed to fetch open information activity records by account.   

The `PreGod` version focuses on a specific wallet address, while the `Node` version expands the functionality to account-based queries with enhanced filtering options.

#### Request Parameters Changes

##### PreGod Version

- `address` (required, string): The wallet address to query.
- `limit` (number, optional): Limits the number of records returned.
- `cursor` (string, optional): Pagination cursor for continuing from a specific point.
- `type` (array of strings, optional): Filters by activity types.
- `tag` (array of strings, optional): Filters by activity tags.
- `network` (array of strings, optional): Filters by network of the open information.
- `platform` (array of strings, optional): Filters by platform of the open information.
- `timestamp` (Time, optional): Filters activities from a specific timestamp.
- `hash` (string, optional): Searches for an activity by hash.
- `hash_list` (array of strings, optional): Filters by a list of activity hashes.
- `include_poap` (boolean, optional): Includes Proof of Attendance Protocols in the results.
- `refresh` (boolean, optional): Forces a refresh of cached data.
- `page` (number, optional): Specifies pagination page.
- `query_status` (boolean, optional): Queries the status of activities.
- `token_id` (string, optional): Filters by specific token ID.
- `count_only` (boolean, optional): Returns only the count of records.
- `action_limit` (number, optional): Limits the number of actions per activity.

##### Node Version

- `account` (required, string): Specifies the open information account to retrieve activities.
- `limit` (integer, optional): Specifies the number of activities to retrieve, with defaults and limits.
- `action_limit` (integer, optional): Specifies the maximum number of actions within an activity.
- `cursor` (string, optional): Used for pagination.
- `since_timestamp` (integer, optional): Retrieves activities starting from this timestamp.
- `until_timestamp` (integer, optional): Retrieves activities up to this timestamp.
- `success` (boolean, optional): Filters activities based on success.
- `direction` (Direction, optional): Filters by the direction of activities.
- `network` (array of Network, optional): Filters by open information network.
- `tag` (array of Tag, optional): Filters by activity tags.
- `type` (array of strings, optional): Filters by activity types.
- `platform` (array of Platform, optional): Filters by platforms involved in the activities.

##### Request Parameters Migration Explanation

In the migration from PreGod API to the Node API, users will encounter several changes to the request parameters, which involve additions, removals, and modifications. Here’s a detailed breakdown of how these parameters have changed:

**Parameters Removed:** 

- `hash`: We provide a new separate API endpoint, /decentralized/tx/{id}, specifically for querying individual activity, removing the need for this functionality in the account activity endpoint.
- `hash_list`: The ability to query multiple hashes simultaneously has been temporarily removed. We plan to introduce a new endpoint in the future that will handle batch hash queries.
- `include_poap`: Special handling for Proof of Attendance Protocols (POAPs) is no longer required in the API, simplifying the data retrieval process.
- `refresh`: Manual refreshes are no longer necessary as the data provided by the API is indexed in real-time by the nodes, ensuring up-to-date information is always available.
- `page`: Pagination is now managed through the cursor parameter, which offers a more efficient and flexible way to navigate through large datasets.
- `query_status`: This parameter has been replaced by the `success` parameter, which directly indicates the success status of activities.
- `token_id`: Support for filtering by token ID has been temporarily removed pending further development.
- `count_only`: Due to the vast amount of open information data, it is challenging to provide accurate count statistics, leading to the removal of this parameter.
- `timestamp`: Timestamp filtering is now handled by `since_timestamp` and `until_timestamp`, allowing for more precise control over the time range of the activities retrieved.

**Parameters Added:**

- `since_timestamp`
- `until_timestamp`
- `success`
- `direction`

**Parameters Modified:**

- `address` -> `account`: Changed terminology to better reflect the open information.
- `limit`, `action_limit`: These parameters were retained but are now more clearly defined with defaults and limits.
- `cursor`: Continues to be used for pagination, emphasizing a cursor-based navigation.
- `type`, `tag`, `network`, `platform`: These remain as filters but have been updated to use schema references for more flexible and robust querying.


These changes streamline the API by removing less frequently used parameters and by refining the functionality of essential parameters, thus enhancing usability and performance. Users need to update their API calls to align with these new parameters, ensuring they leverage the more efficient and flexible Node API for their open information activity queries.


#### Response Changes

##### PreGod Version

Certainly! Here is the updated description of the PreGod version's response structure using Markdown with backticks to highlight the parameters:

### PreGod Version Response Structure

The response structure in the PreGod API provides detailed data associated with activities linked to a specific wallet address. Here’s the structure formatted with backticks to highlight parameters:

- `address_status` (array of string): Contains the refresh status of each queried address, reflecting the current state of the cache for that address. This helps users understand the timeliness and reliability of the data presented.
- `cursor` (string): A string used for pagination, indicating the position from which to continue fetching data in subsequent requests.
- `message` (string): A general message about the response, typically regarding the status of the request.
- `total` (number, nullable): The total number of items that match the query, which may be absent if not applicable.
- `result` (array of `Transaction`): An array of activities relevant to the queried address.
  - `actions` (array of `Action`): Details of actions involved in the activity.
  - `address_from` (string): The originating address of the activity.
  - `address_to` (string): The destination address of the activity.
  - `created_at`, `updated_at`, `timestamp` (`Time`): Timestamps indicating the creation, last update, and the actual time of the activity.
  - `fee` (nullable): The activity fee, which may be null.
  - `hash` (string): The hash of the activity.
  - `network` (string): The open information network on which the activity occurred.
  - `owner` (string): The owner of the address involved in the activity.
  - `platform` (string): The platform associated with the activity.
  - `success` (boolean, nullable): Indicates whether the activity was successful.
  - `tag` (string): A tag related to the activity.
  - `type` (string): The type of the activity.

##### Node Version:

In the Node API, the response is structured to focus more on activities rather than individual activities, reflecting a broader view of account interactions:


- `data` (array of Activity): Lists the activities associated with the account.
  - `actions` (array): Lists actions within the activity.
  - `calldata` (Calldata): Details of the call made during the activity.
  - `direction` (Direction): The direction of the activity.
  - `fee` (Fee): Detailed fee information for the activity.
  - `from`, `to` (string): The originating and destination addresses of the activity.
  - `id` (string): A unique identifier for the activity.
  - `index` (integer): The index position of the activity in the list.
  - `network` (Network): The network on which the activity occurred.
  - `owner` (string): The owner of the account involved in the activity.
  - `platform` (Platform): The platform related to the activity.
  - `success` (boolean): Indicates the success of the activity.
  - `tag` (Tag): A tag categorizing the activity.
  - `type` (string): The type of activity.   
  - `timestamp` (integer): The timestamp when the activity occurred.
  - `total_actions` (integer): The total number of actions within the activity.
- `meta` (MetaCursor): Contains metadata such as the pagination cursor for continuing the data fetch.
  - `cursor` (string): The cursor for the next set of results.

##### Response Migration Explanation


**Fields Removed:**
- `message`: Removed because the Node version does not currently have any special informational messages to return.
- `total`: Not provided in the Node version due to the large volume of open information data, making count functionality impractical.
- `address_status`: No longer needed as the Node version does not use caching for address indexing, unlike the PreGod version.

**Fields Added:**

- `calldata`
- `direction`
- `id`
- `index`
- `total_actions`


**Fields Modified:**

- `hash` -> `id`
- `address_from` -> `from`, `address_to` -> `to`
- `platform`, `network`, `tag`, `type`: While these fields exist in both versions, in the Node API they are possibly enhanced with richer schema definitions or standardized formats to ensure consistency across different parts of the API.
- `timestamp`: Changed from a potentially more complex `Time` object in the PreGod version to a simpler integer format in the Node version.
