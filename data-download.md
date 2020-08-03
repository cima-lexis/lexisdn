# Kind of dataset needed 



## GFS data

Various components of the simulations requires GFS data as input. They can be downloaded using a web service provided by noaa.

GFS are organized in datasets, with every dataset containing 384 hours of forecast, 1 hour of initial conditions and 1 file of analysis (available at 00, 06, 12 and 18UTC a given day).

Also, the nooa service to download the data is able to filter a subregion of interested, thus greatly reducing size of downloaded data.

Each file can be downloaded using following URLs:

```curl
https://nomads.ncep.noaa.gov/cgi-bin/filter_gfs_0p25.pl?file=gfs.t${ds_hour}z.pgrb2.0p25.${fileno}&all_lev=on&all_var=on&subregion=&leftlon=${leftlon}&rightlon=${rightlon}&toplat=${toplat}&bottomlat=${bottomlat}&dir=%2Fgfs.${ds_date}%2F${ds_hour}
```

where the various variable represents: 

- **ds_date**: the date of start of the dataset (format YYYYMMDD)
- **ds_hour**: the hour of start of the dataset (format HH, eg. 00, 06, 12, 18)
- **leftlon,rightlon,toplat,bottomlat**: boundaries of requested subregion
- **fileno**: file of the dataset to download, set to:
     _anl_ for analysis file
     _f000_ for initial conditions
     _f000_-_f384_ for forecasted hours

from now on, when we speak to download GFS with a **ds_date**, **ds_hour** and **fileno** we mean to download the GFS file with the URL built using that parameters and the domain data needed.



## Wunderground personal weather stations data

These datasets are available on CIMA servers, and they can be downloaded using CIMA WEBDROPS service. They contains hourly observations data from a network of weather stations provided by IBM Weather Company.

Dataset is organized in multiple "sensors" classes. For each sensor class needed, a registry of all available stations can be retried. Data of observations is downloaded issuing a separate HTTP call for each of the sensors needed.

### Registry download

 You can download the registry for a sensor by issuing a GET request (upon [authentication](#WEBDROPS-Authentication )) to following URL:

```curl 
http://webdrops.cimafoundation.org/app/sensors/list/<sensor_class>
```

where **sensor_class** is the name of the sensor class  registry you want to read. 



The registry contains informations of all available stations for specified sensor class.



#### Example response body:

```json
[
  {
    "id": "-1937152789_2",
    "name": "Giardino Botanico Celle",
    "lat": 44.343433,
    "lng": 8.54158,
    "mu": "mm"
  },
  {
    "id": "-1937157087_2",
    "name": "Località Beo",
    "lat": 44.05301,
    "lng": 8.088548,
    "mu": "mm"
  },
  {
    "id": "-1937156901_2",
    "name": "Suvero",
    "lat": 44.265083,
    "lng": 9.776026,
    "mu": "mm"
  },

...

]
```





### Punctual observations download

You can download observations data for a list of sensors  by issuing (upon [authentication](#WEBDROPS-Authentication )) a **POST** request to following URL: 

```curl
http://webdrops.cimafoundation.org/app/sensors/data/<sensor_class>/?from=<from>&to=<to>&aggr=<aggr>
```

where the variables used are:

* **sensor_class**: name of the class of sensors  to download
* **from**: start of the temporal range of observations to download [format **YYYMMDDHHmm**]
* **to**: end of the temporal range of observations to download [format **YYYMMDDHHmm**]
* **aggr**: aggregation granularity in seconds. Observations data will be aggregated with this timespan

The body must be in JSON format, and must contains these additional parameters:

- **sensors**: array of sensors id to download. They can be obtained downloading the [sensors registry](Registry-download), or provided otherwise.

  

#### Example request body:

```json
{
	"sensors": ["7272_2","51243_1","27613_2"],
}
```

#### Example response body:

```json
[
    {
        "sensorId": "479098258_2",
        "timeline": [
            "202004160001",
            "202004160101"
        ],
        "values": [
            0.0,
            0.0
        ]
    },
    {
        "sensorId": "479124472_2",
        "timeline": [
            "202004160001",
            "202004160101"
        ],
        "values": [
            0.0,
            0.0
        ]
    },
    {
        "sensorId": "479098260_2",
        "timeline": [],
        "values": []
    },
    {
        "sensorId": "479124476_2",
        "timeline": [
            "202004160001",
            "202004160101"
        ],
        "values": [
            0.0,
            0.0
        ]
    }
]
```



### Interpolated maps

These datasets are available on CIMA servers, and they can be downloaded using CIMA WEBDROPS service. They contains weather stations observations interpolated in a map in NETCDF format.

They can be downloaded  by issuing (upon [authentication](#WEBDROPS-Authentication )) a **GET** request to following URL: 

```curl
http://webdrops.cimafoundation.org/app/sensors/map/<sensor_class>/?from=<from>&to=<to>
```

where the variables used are:

* **sensor_class**: name of the class of sensors  to download
* **from**: start of the temporal range of observations to download [format **YYYMMDDHHmm**]
* **to**: end of the temporal range of observations to download [format **YYYMMDDHHmm**]

The HTTP response will be a byte streams containing an interpolated map of the observations. The format of returned file is NETCDF.



## Radar data

These datasets are available on CIMA servers, and they can be downloaded using CIMA WEBDROPS service. They contains radar data over Italy and France.



### Radar timeline

You can download the timeline for a radar by issuing a GET request (upon [authentication](#WEBDROPS-Authentication )) to following URL:

```curl 
http://webdrops.cimafoundation.org/app/coverages/<radar_dataset>/?from=<from>&to=<to>
```

where the variables used are: 

* **radar_dataset** is the name of radar you want to read. 
* **from**: start of the temporal range of radar to include [format **YYYMMDDHHmm**]
* **to**: end of the temporal range of radar to include [format **YYYMMDDHHmm**]

The radar timeline contains date/time instants of all radar datasets present. This is useful to download only the radar data nearer to a needed instant in time.

#### Example response body:

```json
[
    "202004160001",
    "202004160101"
]
```



### Radar data download

You can download data for radars by issuing a GET request (upon [authentication](#WEBDROPS-Authentication )) to following URL:

```curl 
http://webdrops.cimafoundation.org/app/coverages/<radar_dataset>/<date>/<varname>/-/all
```

where the variables used are: 

* **radar_dataset** is the name of radars you want to read. 
* **date**: instant in time of radar to download [format **YYYMMDDHHmm**]
* **varname**: name of the variable to include in the downloaded NETCDF file.



The HTTP response will be a byte streams containing an interpolated map of the observations. The format of returned file is NETCDF.

## 





## WEBDROPS Authentication 

`WEBDROP` use Key Kloak for authentication. You need to obtain a token in order to authenticate any following requests.

To get a token, issue a POST request to `https://testauth.cimafoundation.org/auth/realms/webdrops/protocol/openid-connect/token`
the request body must be formatted as `application/x-www-form-urlencoded` and must include the relative Content-Type header.

The body must include following fields:

* grant_type=password
* username=[provided to you by CIMA]
* password=[provided to you by CIMA]
* client_id=webdrops

On success, you'll receive a JSON formatted response, containing an "access_token" field you can use to authenticate subsequent requests, using a Bearer Token authorization header.

### Example request body:

```
grant_type:password
username:andrea.parodi@cimafoundation.org
password:nottherealpassword
client_id:webdrops
```

### Example response body:

```json
{
    "access_token": "fdgfdgfdgfdg ...",
    "expires_in": 300,
    "refresh_expires_in": 1800,
    "refresh_token": "dsfdfdsfs ...",
    "token_type": "bearer",
    "not-before-policy": 0,
    "session_state": "74de0e22-7ff5-42d3-82de-e0821c52aeac",
    "scope": "profile email"
}
```







# Dataset needed by component



## WPS docker container



WPS need the following external data to run:

* Static geographic data - https://www2.mmm.ucar.edu/wrf/users/download/get_sources_wps_geog.html
* GFS model data - used for italian case studies
* IFS model data - used for french case studies

For usage with WRF, we need N hours of GFS or IFS data, starting from the beginning of the run,
where N is the number of forecast hours we want to produce.

For WRF with data assimilation, we need in addition the data for 3 and 6 hours before start of WRF

>  _keep in mind that for simulation involving risico, there will be other GFS or IFS files needed. See the [chapter on risico](#Risico) for further details_



### Examples

To forecast from 2020-06-10 00:00 upto 2020-06-11 00:00, we'll need these GFS data:

* 2020-06-09 18:00	-	needed only for simulations with data assimilation  

* 2020-06-09 21:00	-	needed only for simulations with data assimilation

* from 2020-06-10 00:00 upto 2020-06-11 00:00

  

In other words, referring to the variables [specified here](#GFS data) we will have to download following GFS URLs:

1) **ds_date** = 20200609, **ds_hour** = 18 ,**fileno** 0

1) **ds_date** = 20200609, **ds_hour** = 18 ,**fileno** 3

2) **ds_date** = 20200609, **ds_hour** = 18 ,**fileno** from 6 to 30

Here you can find an example bash script that download all needed GFS files for a WRF run:

https://github.com/cima-lexis/wps.docker/blob/download-gfs/wps.gfs/gfs-download.sh



## WRFDA

WRFDA can assimilate data of wunderground stations. These can be downloaded issuing various requests as explained below.
All this requests need to be authenticate using a Bearer Token authorization header obtained as explained above.
Service base URL is `http://webdrops.cimafoundation.org/app/`. 

### Retrieve wunderground station data

WRFDA can assimilate the following variables from wunderground stations:

* IGROMETRO -> relative humidity
* TERMOMETRO -> air temperature
* PLUVIOMETRO -> precipitation
* ANEMOMETRO -> wind speed
* BAROMETRO -> pressure
* DIREZIONEVENTO -> wind direction

```<Dewetra sensor class> -> <variable name>```

As explained [in the relative chapter](#Registry download) to download a set of observations we first need to download a registry of sensors to obtain the IDs of the sensors we want.

However, for all the case dates identified for the LEXIS projects, we will provide a list of IDs to download, that we prepared in CIMA server. We pre-filtered the stations to exclude the ones that are not performing well for the date of the study. 

Alternatively if you have to assimilate a different date, you can download a list of all available stations by downloading the sensors registry as explained. 

> In a run that start at moment N, WRFDA need to assimilate 3 instant in time: N, N-3H, N-6H.
> These 3 moments can be downloaded with 3 sets of http requests for the 3 specific moment, or otherwise with a unique http call that download all data from N-6H to N.

Sensors data must be downloaded specified an aggregation time of 60 seconds (1m)

[here](https://github.com/cima-lexis/lexisdn/blob/master/fetcher/wrfda-sensors.go) you can find an example script that downloads all wunderground files needed.



### Retrieve radar data

Our WRFDA run also assimilates radar data, with the same timing as the personal weather stations (N, N-3H, N-6H).
For each one of these needed instants, three sets of variables are needed for the different levels. These variables are CAPPI2, CAPPI3, CAPPI5.

You can download the datasets as [explained here](#Radar data)

[here](https://github.com/cima-lexis/lexisdn/blob/master/fetcher/wrfda-radars.go) you can find an example script that downloads all radar files needed.



## Continuum

Continuum needs wunderground station data from 60H before the start of the run. The datasets can be downloaded as [explained here](#Puntual observations download).

Needed sensors classes are:

* IGROMETRO    -> relative humidity

* TERMOMETRO -> air temperature

* PLUVIOMETRO -> precipitation

* ANEMOMETRO -> wind speed

* RADIOMETRO -> incoming radiation


As for WRFDA, we will provide a file with IDs of sensors to download.

Sensors data must be downloaded with an aggregation time of 3600 seconds (1H)

[here](https://github.com/cima-lexis/lexisdn/blob/master/fetcher/continuum.go) you can find an example script that downloads all files needed by Continuum simulation



## Risico

Risico could make use of the following wunderground sensors classes:

* TERMOMETRO -> air temperature
* PLUVIOMETRO -> precipitation
* IGROMETRO    -> relative humidity

Risico needs sensors data for 72H before the start of the run.
Data must be downloaded as interpolated maps, as [explained here](#Interpolated maps)

[here](https://github.com/cima-lexis/lexisdn/blob/master/fetcher/risico.go) you can find an example script that downloads all files needed by Risico simulation






