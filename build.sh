source /etc/profile.d/modules.sh
module load netcdf-hdf5-all/4.7_hdf5-1.10-gcc8-serial
module unload intel-mkl/2019 intel-mpi/2019-intel intel/19.0.5  
go build -o lexisdn ./cli
