

# Bartend


## <a name=""></a>/
Bartender is an app that will mix a drink by turning on a set of 
liquid pumps


  
* **[pump[…]](#/pump)** - direct access to hardware pumps. 

  
* **[recipe[…]](#/recipe)** - available list of drinks known to this bartender. 

  
* **liquids** `string[]` - all the unique liquids found across all recipes. 

  
* **[current](#/current)** - current drink recipe in progress. 







## <a name="/pump"></a>/pump={id}/
direct access to hardware pumps


  
* **id** `int32` - ids from 0-8 for the installed pumps. 

  
* **gpioPin** `int32` - Raspberry PI GPIO pin number. Do not change unless you&#39;re rewired your Pi.. 

  
* **timeToVolumeRatioMs** `int32` - number of millisecs to turn on pump to pour one milliliter. 

  
* **liquid** `string` - name of the liquid. If name does not match (case sensitive) recipe 
ingredient exactly, it will not recognized as same liquid.. 



### Actions:

* <a name="/pump/on"></a>**/pump={id}/on** - turn on the pump
 
  


  


* <a name="/pump/off"></a>**/pump={id}/off** - turn off the pump
 
  


  







## <a name="/recipe"></a>/recipe={id}/
available list of drinks known to this bartender


  
* **id** `int32` - unique id for this drink. 

  
* **name** `string` - name of drink. 

  
* **description** `string` - description of drink. 

  
* **[ingredient[…]](#/recipe/ingredient)** - ingredients for drink. 

  
* **madeCount** `int32` - number of times the this recipee was made. 



### Actions:

* <a name="/recipe/make"></a>**/recipe={id}/make** - make this drink
 
  
#### Input:
> * **multiplier** `decimal64` - Optional, make it a double by passing 2 or sample it by passing 
0.1


  







## <a name="/recipe/ingredient"></a>/recipe={id}/ingredient={liquid}/
ingredients for drink


  
* **liquid** `string` - liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid.. 

  
* **amount** `decimal64` - in ounces of liquid. 

  
* **weight** `int32` - in grams of liquid determined by amount and estimate of grams per ounce. 







## <a name="/current"></a>/current/
current drink recipe in progress


  
* **percentComplete** `int32` - 0-100 percent complete with drink. 

  
* **complete** `boolean` - true if drink has finished being made. 

  
* **[auto[…]](#/current/auto)** - ingredients from recipe that are available in the pumps and would
be poured automatically. 

  
* **[manual[…]](#/current/manual)** - ingredients from recipe that are not available in the pumps and would
need to be poured manually by user. 



### Actions:

* <a name="/current/stop"></a>**/current/stop** - Stop making drink where ever in the process it is
 
  


  





### Events:

* <a name="/current/update"></a>**/current/update** - 

 	
> * **percentComplete** `int32` - 0-100 percent of how far along is this drink done making	
> * **complete** `boolean` - is this drink ready for drinking
> * **manual[…]** - ingredients that should be poured manually
>     * **id** - unique id for this ingredient 
>     * **complete** - is this ingredient done pouring 
> * **auto[…]** - ingredients that should be poured automatically
>     * **id** - unique id for this ingredient 
>     * **complete** - is this ingredient done pouring 
>     * **percentComplete** - 0-100 percent of how far along is this ingredient is done pouring 





## <a name="/current/auto"></a>/current/auto={id}/
ingredients from recipe that are available in the pumps and would
be poured automatically


  
* **id** `int32` - unique id for this ingredient. 

  
* **[ingredient](#/current/auto/ingredient)** - liquid ingredient. 

  
* **complete** `boolean` - is this liquid done pouring. 

  
* **percentComplete** `int32` - what percentage is this liquid complete pouring. 







## <a name="/current/auto/ingredient"></a>/current/auto={id}/ingredient/
liquid ingredient


  
* **liquid** `string` - liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid.. 

  
* **amount** `decimal64` - in ounces of liquid. 

  
* **weight** `int32` - in grams of liquid determined by amount and estimate of grams per ounce. 







## <a name="/current/manual"></a>/current/manual={id}/
ingredients from recipe that are not available in the pumps and would
need to be poured manually by user


  
* **id** `int32` - unique id for this ingredient. 

  
* **[ingredient](#/current/manual/ingredient)** - liquid ingredient. 

  
* **complete** `boolean` - is this liquid done pouring. 



### Actions:

* <a name="/current/manual/done"></a>**/current/manual={id}/done** - mark step done
 
  


  







## <a name="/current/manual/ingredient"></a>/current/manual={id}/ingredient/
liquid ingredient


  
* **liquid** `string` - liquid name. If name does not match (case sensitive) pump 
liquid exactly, it will not recognized as same liquid.. 

  
* **amount** `decimal64` - in ounces of liquid. 

  
* **weight** `int32` - in grams of liquid determined by amount and estimate of grams per ounce. 







