# Week9-assignment

Created the TrimmedMean package, and this repo importing/running it.

To import TrimmedMean package, add to go.mod:  
require github.com/sdnkearns/TrimmedMean v0.1.0  

To use TrimmedMean:  
trimmedmean.TrimmedMean(T, lowTrim, highTrim)  

Inputs:  
 -number slice T (accepts various numerical variable types)  
 -lowTrim (percent of lowest values in T to be trimmed)  
 -highTrim (percent of highest values in T to be trimmed. If no value provided, will use the same value as lowTrim)  

outputs:  
 -float64 value containing the mean of slice T after trimming lowTrim% of the lowest values and highTrim% of the highest values in T

 run main.go as:  
 go run main.go

 build main.go as:  
 go build main.go

The The R calculation for trimmed means was slightly more accurate, being only 0.11 away from the target mean of 100, where the Go version was 0.24 away, however go was noticably faster when increasing the size of the slice to around 10000.  
The difference in accuracy could also be the result of differences in the random number generation. I used the mersenne twister algorithm used in assignment 8 for the random number generation, but that may not guarantee exact results.

I used a claude to generate the unit tests in my TrimmedMean repository, as well as to help replicate the random number generation the R code uses in go.
