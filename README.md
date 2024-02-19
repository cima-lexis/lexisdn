# lexis-docs

Italian radar & wunderground weather stations to WRF

This module can be used to download and convert wunderground weather stations observations and italian radar data into ascii WRF format.

## Installation
The module use go-netcdf to read a netcdf containing world orography data.

In order to use it, you need the developer version of the library provided by your distribution installed.

On ubuntu you can install it with:

sudo apt install libnetcdf-dev
On Typhoon, it can be loaded with the WRF-KIT2 module:

module load gcc-8.3.1/WRF-KIT2

The orography data is used to calculate the altitude of every weather station given their latitude and longitude coordinates.

You can download the orography file from https://zenodo.org/record/4607436/files/orog.nc

The file must be saved in path ~/.dewetra2wrf/orog.nc

## Usage on CIMA Typhoon
An orography file is already usable by wrfprod user: /data/safe/home/wrfprod/.dewetra2wrf/orog.nc.

lexisdn is already present in /data/safe/home/wrfprod/bin/lexisdn

## Command line usage
This module implements a console command that can be used to convert observations and radars to ascii WRF format.

Usage of lexisdn:

Usage: lexisdn STARTDATE [DOWNLOAD_TYPE ...]
	STARTDATE - Satrt date/time of the simulation, in format YYYYMMDDHH
	DOWNLOAD_TYPE - types of data to download. One of "WRFIT" | "WRFITDPC" | "WRFFR"

This commands require following environment variable to be set:
  WEBDROPS_USER			-	webdrops user
  WEBDROPS_PWD			-	webdrops password
  WEBDROPS_CLIENT_ID	-	webdrops client id
  WEBDROPS_AUTH_URL		-	URL for KeyCloak authentication
  WEBDROPS_URL			-	base URL for all webdrops endpoints

