
# Bartend

Bartender is an app that will mix a drink by turning on a set of 
liquid pumps

<details><summary>API Usage Notes:</summary>

#### General API Usage Notes
* `DELETE` implementation may be disallowed or ignored depending on the context
* Lists use `../path={key}/...` instead of `.../path/key/...` to avoid API name collision

#### `GET` Query Parameters

These parameters can be combined.

> | param                            | description | example |
> |----------------------------------|-------------|---------|
> | `content=[non-config\|config]` | Show only read-only fields or only read/write fields |   `.../path?content=config`|
> | `fields=field1;field2` | Return a portion of the data limited to fields listed | `.../path?fields=user%2faddress` |
> | `depth=n` | Return a portion of the data limited to depth of the hierarchy | `.../path?depth=1`
> | `fc.xfields=field1;fields` | Return a portion of the data excluding the fields listed | `.../path?fc.xfields=user%2faddress` |
> | `fc.range=field!{startRow}-[{endRow}]` | For lists, return only limited number of rows | `.../path?fc.range=user!10-20` 

</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:bartend</b></code> Bartender is an app that will mix a drink by turning on a set of 
liquid pumps</summary>

#### bartend


**GET Response Data**
````json
{
  "pump":[{
     "id":0,
     "gpioPin":0,
     "timeToVolumeRatioMs":0,
     "liquid":""
  }],
  "available":[{
     "name":"",
     "description":"",
     "ingredient":[{
        "liquid":"",
        "amount":0
     }]
  }],
  "recipe":[{
     "name":"",
     "description":"",
     "ingredient":[{
        "liquid":"",
        "amount":0
     }],
     "madeCount":0
  }],
  "liquids":["", "..."],
  "drink":{
  }}
````

**PUT, POST Request Data**
````json
{
  "pump":[{
     "id":0,
     "gpioPin":0,
     "timeToVolumeRatioMs":0,
     "liquid":""
  }],
  "recipe":[{
     "name":"",
     "description":"",
     "ingredient":[{
        "liquid":"",
        "amount":0
     }],
     "madeCount":0
  }]}
````

**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | pump.id | int32  |  ids from 0-8 for the installed pumps |  |
> | pump.gpioPin | int32  |  Raspberry PI GPIO pin number. Do not change unless you&#39;re rewired your Pi. |  |
> | pump.timeToVolumeRatioMs | int32  |  number of millisecs to turn on pump to pour one milliliter |  |
> | pump.liquid | string  |  name of the liquid. If name does not match (case sensitive) recipe 
ingredient exactly, it will not recognized as same liquid. |  |
> | available.name | string  |  name of drink | r/o |
> | available.description | string  |  description of drink | r/o |
> | available.ingredient.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | available.ingredient.amount | decimal64  |  in ounces of liquid | r/o |
> | recipe.name | string  |  name of drink |  |
> | recipe.description | string  |  description of drink |  |
> | recipe.ingredient.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | recipe.ingredient.amount | decimal64  |  in ounces of liquid | r/o |
> | recipe.madeCount | int32  |  number of times the this recipe was made | r/o |
> | liquids | string[]  |  all the unique liquids found across all recipes | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:bartend

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:bartend

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:bartend

# delete current data
curl -X DELETE https://server/restconf/data/bartend:bartend
````
</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:pump</b></code> direct access to hardware pumps</summary>

#### pump

**GET Response Data / PUT, POST Request Data**
````json
{"pump":[{ 
  "id":0,
  "gpioPin":0,
  "timeToVolumeRatioMs":0,
  "liquid":""}, {"..."}]}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | id | int32  |  ids from 0-8 for the installed pumps |  |
> | gpioPin | int32  |  Raspberry PI GPIO pin number. Do not change unless you&#39;re rewired your Pi. |  |
> | timeToVolumeRatioMs | int32  |  number of millisecs to turn on pump to pour one milliliter |  |
> | liquid | string  |  name of the liquid. If name does not match (case sensitive) recipe 
ingredient exactly, it will not recognized as same liquid. |  |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:pump

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:pump

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:pump

# delete current data
curl -X DELETE https://server/restconf/data/bartend:pump
````
</details>




<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:pump={id}</b></code> direct access to hardware pumps</summary>

#### pump={id}

**GET Response Data / PUT, POST Request Data**
````json
{
  "id":0,
  "gpioPin":0,
  "timeToVolumeRatioMs":0,
  "liquid":""}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | id | int32  |  ids from 0-8 for the installed pumps |  |
> | gpioPin | int32  |  Raspberry PI GPIO pin number. Do not change unless you&#39;re rewired your Pi. |  |
> | timeToVolumeRatioMs | int32  |  number of millisecs to turn on pump to pour one milliliter |  |
> | liquid | string  |  name of the liquid. If name does not match (case sensitive) recipe 
ingredient exactly, it will not recognized as same liquid. |  |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:pump={id}

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:pump={id}

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:pump={id}

# delete current data
curl -X DELETE https://server/restconf/data/bartend:pump={id}
````
</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:available</b></code> list of only the drinks that are available to make given the available ingredients at the pumps</summary>

#### available


**GET Response Data**
````json
{"available":[{ 
  "name":"",
  "description":"",
  "ingredient":[{
     "liquid":"",
     "amount":0
  }]}, {"..."}]}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | name | string  |  name of drink | r/o |
> | description | string  |  description of drink | r/o |
> | ingredient.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | ingredient.amount | decimal64  |  in ounces of liquid | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:available

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:available

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:available

# delete current data
curl -X DELETE https://server/restconf/data/bartend:available
````
</details>




<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:available={name}</b></code> list of only the drinks that are available to make given the available ingredients at the pumps</summary>

#### available={name}


**GET Response Data**
````json
{
  "name":"",
  "description":"",
  "ingredient":[{
     "liquid":"",
     "amount":0
  }]}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | name | string  |  name of drink | r/o |
> | description | string  |  description of drink | r/o |
> | ingredient.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | ingredient.amount | decimal64  |  in ounces of liquid | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:available={name}

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:available={name}

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:available={name}

# delete current data
curl -X DELETE https://server/restconf/data/bartend:available={name}
````
</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:available={name}/ingredient</b></code> ingredients for drink</summary>

#### available={name}/ingredient


**GET Response Data**
````json
{"ingredient":[{ 
  "liquid":"",
  "amount":0}, {"..."}]}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | amount | decimal64  |  in ounces of liquid | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:available={name}/ingredient

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:available={name}/ingredient

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:available={name}/ingredient

# delete current data
curl -X DELETE https://server/restconf/data/bartend:available={name}/ingredient
````
</details>




<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:available={name}/ingredient={liquid}</b></code> ingredients for drink</summary>

#### available={name}/ingredient={liquid}


**GET Response Data**
````json
{
  "liquid":"",
  "amount":0}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | amount | decimal64  |  in ounces of liquid | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:available={name}/ingredient={liquid}

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:available={name}/ingredient={liquid}

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:available={name}/ingredient={liquid}

# delete current data
curl -X DELETE https://server/restconf/data/bartend:available={name}/ingredient={liquid}
````
</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:recipe</b></code> list of all drinks known to this bartender</summary>

#### recipe


**GET Response Data**
````json
{"recipe":[{ 
  "name":"",
  "description":"",
  "ingredient":[{
     "liquid":"",
     "amount":0
  }],
  "madeCount":0}, {"..."}]}
````

**PUT, POST Request Data**
````json
{"recipe":[{ 
  "name":"",
  "description":"",
  "ingredient":[{
     "liquid":"",
     "amount":0
  }]},{"..."}]}
````

**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | name | string  |  name of drink |  |
> | description | string  |  description of drink |  |
> | ingredient.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | ingredient.amount | decimal64  |  in ounces of liquid | r/o |
> | madeCount | int32  |  number of times the this recipe was made | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:recipe

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:recipe

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:recipe

# delete current data
curl -X DELETE https://server/restconf/data/bartend:recipe
````
</details>




<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:recipe={name}</b></code> list of all drinks known to this bartender</summary>

#### recipe={name}


**GET Response Data**
````json
{
  "name":"",
  "description":"",
  "ingredient":[{
     "liquid":"",
     "amount":0
  }],
  "madeCount":0}
````

**PUT, POST Request Data**
````json
{
  "name":"",
  "description":"",
  "ingredient":[{
     "liquid":"",
     "amount":0
  }]}
````

**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | name | string  |  name of drink |  |
> | description | string  |  description of drink |  |
> | ingredient.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | ingredient.amount | decimal64  |  in ounces of liquid | r/o |
> | madeCount | int32  |  number of times the this recipe was made | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:recipe={name}

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:recipe={name}

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:recipe={name}

# delete current data
curl -X DELETE https://server/restconf/data/bartend:recipe={name}
````
</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:recipe={name}/ingredient</b></code> ingredients for drink</summary>

#### recipe={name}/ingredient


**GET Response Data**
````json
{"ingredient":[{ 
  "liquid":"",
  "amount":0}, {"..."}]}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | amount | decimal64  |  in ounces of liquid | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:recipe={name}/ingredient

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:recipe={name}/ingredient

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:recipe={name}/ingredient

# delete current data
curl -X DELETE https://server/restconf/data/bartend:recipe={name}/ingredient
````
</details>




<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:recipe={name}/ingredient={liquid}</b></code> ingredients for drink</summary>

#### recipe={name}/ingredient={liquid}


**GET Response Data**
````json
{
  "liquid":"",
  "amount":0}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. | r/o |
> | amount | decimal64  |  in ounces of liquid | r/o |

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:recipe={name}/ingredient={liquid}

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:recipe={name}/ingredient={liquid}

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:recipe={name}/ingredient={liquid}

# delete current data
curl -X DELETE https://server/restconf/data/bartend:recipe={name}/ingredient={liquid}
````
</details>





<details>
 <summary><code>[GET|PUT|POST|DELETE]</code> <code><b>restconf/data/bartend:drink</b></code> current drink recipe in progress</summary>

#### drink

**GET Response Data / PUT, POST Request Data**
````json
{}
````



**Data Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|

**Responses**
> | http method  |  request body  | response body |
> |--------------|----------------|---------------|
> | `POST`       |  *JSON data*   | - none -      |
> | `PUT`       |  *JSON data*   | - none -      |
> | `GET`       |  - none -      | *JSON data*   |
> | `DELETE`     |  - none -      | - none -      |

**HTTP response codes**
> | http code |  reason for code    |
> |-----------|---------------------|
> | 200       | success             |
> | 401       | not authorized      |
> | 400       | invalid request     |
> | 404       | data does not exist |
> | 500       | internal error      |

**Examples**
````bash
# retrieve data
curl https://server/restconf/data/bartend:drink

# update existing data
curl -X PUT -d @data.json https://server/restconf/data/bartend:drink

# create new data
curl -X POST -d @data.json https://server/restconf/data/bartend:drink

# delete current data
curl -X DELETE https://server/restconf/data/bartend:drink
````
</details>




  <details>
 <summary><code>[POST]</code> <code><b>restconf/data/bartend:pump={id}/on</b></code> turn on the pump</summary>
 
#### pump={id}/on

 **Request Body**
    
  *none*
    

**Response Body**
    
  *none*
    

**HTTP response codes**

> | http code |  reason for code |
> |-----------|------------------|
> | 200       | success          |
> | 401       | not authorized   |
> | 400       | invalid request  |
> | 404       | data does not exist |
> | 500       | internal error   |

**Examples**
````bash
# call function
curl -X POST  https://server/restconf/data/bartend:pump={id}/on
````
  </details>

  <details>
 <summary><code>[POST]</code> <code><b>restconf/data/bartend:pump={id}/off</b></code> turn off the pump</summary>
 
#### pump={id}/off

 **Request Body**
    
  *none*
    

**Response Body**
    
  *none*
    

**HTTP response codes**

> | http code |  reason for code |
> |-----------|------------------|
> | 200       | success          |
> | 401       | not authorized   |
> | 400       | invalid request  |
> | 404       | data does not exist |
> | 500       | internal error   |

**Examples**
````bash
# call function
curl -X POST  https://server/restconf/data/bartend:pump={id}/off
````
  </details>

  <details>
 <summary><code>[POST]</code> <code><b>restconf/data/bartend:available={name}/make</b></code> make this drink</summary>
 
#### available={name}/make

 **Request Body**
    
      
````json
{
  "multiplier":0
}
````

**Request Body Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | multiplier | decimal64  |  Optional, make it a double by passing 2 or sample it by passing 
0.1 |  |
    

**Response Body**
    
  *none*
    

**HTTP response codes**

> | http code |  reason for code |
> |-----------|------------------|
> | 200       | success          |
> | 401       | not authorized   |
> | 400       | invalid request  |
> | 404       | data does not exist |
> | 500       | internal error   |

**Examples**
````bash
# call function
curl -X POST -d @request.json] https://server/restconf/data/bartend:available={name}/make
````
  </details>

  <details>
 <summary><code>[POST]</code> <code><b>restconf/data/bartend:drink/stop</b></code> Stop making drink where ever in the process it is</summary>
 
#### drink/stop

 **Request Body**
    
  *none*
    

**Response Body**
    
  *none*
    

**HTTP response codes**

> | http code |  reason for code |
> |-----------|------------------|
> | 200       | success          |
> | 401       | not authorized   |
> | 400       | invalid request  |
> | 404       | data does not exist |
> | 500       | internal error   |

**Examples**
````bash
# call function
curl -X POST  https://server/restconf/data/bartend:drink/stop
````
  </details>

  



  <details>
 <summary><code>[GET]</code> <code><b>restconf/data/bartend:drink/update</b></code> status update for the drink in progress. Can be a firehose so use
notification filters as defined in IETF protocols if you limit the amount</summary>

#### drink/update

**Response Stream** [SSE Format](https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events)

````
data: {first JSON message all on one line followed by 2 CRLFs}

data: {next JSON message with same format all on one line ...}
````

Each JSON message would have following data
````json
{
  "name":"",
  "percentComplete":0,
  "complete":false,
  "aborted":false,
  "pour":[{
     "pumpId":0,
     "liquid":"",
     "amount":0,
     "complete":false,
     "percentComplete":0
  }],
  "pumpId":0,
  "liquid":"",
  "amount":0,
  "complete":false,
  "percentComplete":0
}
````

**Response Body Details**

> | field   |  type  |  Description |  Details |
> |---------|--------|--------------|----------|
> | name | string  |  recipe name |  |
> | percentComplete | int32  |  0-100 percent complete with drink |  |
> | complete | boolean  |  true if drink has finished being made |  |
> | aborted | boolean  |  was the drink stopped mid pour |  |
> | pour.pumpId | int32  |  unique id for this tap |  |
> | pour.liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. |  |
> | pour.amount | decimal64  |  in ounces of liquid |  |
> | pour.complete | boolean  |  is this liquid done pouring |  |
> | pour.percentComplete | int32  |  what percentage is this liquid complete pouring |  |
> | pumpId | int32  |  unique id for this tap |  |
> | liquid | string  |  liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid. |  |
> | amount | decimal64  |  in ounces of liquid |  |
> | complete | boolean  |  is this liquid done pouring |  |
> | percentComplete | int32  |  what percentage is this liquid complete pouring |  |

**Example**
````bash
# retrieve data stream, adjust timeout for slower streams
curl -N https://server/restconf/data/bartend:drink/update
````

</details>
  

