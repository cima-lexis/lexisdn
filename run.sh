#!/bin/bash
source /etc/profile.d/modules.sh
module load netcdf-hdf5-all/4.7_hdf5-1.10-gcc8-serial
module unload intel-mpi/2019.8.254
module unload intel-mpi/2019-intel
module unload intel-mkl/2019 
module unload intel/19.0.5  
~/repos/lexisdn/lexisdn $@
