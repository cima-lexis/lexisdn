# design

`lexisdn` is composed by three main
modules: 

* **webdrops** abstracts low level HTTP interaction with webdrops server.
* **fetcher** using abstractions provided by `webdrops`, fetcher module orchestrate fetching of all datasets required by various kind of simulation:
WrfdaRadars, ContinuumSensors, RisicoSensorsMaps, WrfdaSensors 

* **conversion** takes care of converting italian radars and wunderground datasets in final wrf ASCII format.