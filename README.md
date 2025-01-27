# app-go-autocomplete-nappi

* an init() function is used to:
  * read 450,000 entries from a text file and hold it in-memory.
  * create an in-memory index of unique words from entries.
* a handler
  * receives requests for autocomplete, which can be one or more words, with each word being at least 3 characters long.
  * uses the index to search for entries that contain the given search words
  * returns the list of matching entries.

## download data

19/01/2025

https://www.medikredit.co.za/products-and-services/nappi/nappi-public-domain-file/


## example usage

curl http://localhost:8080/autocomplete?term=left

curl http://localhost:8080/autocomplete?term=clacee -o "out.json"

curl http://localhost:8080/autocomplete?term=saturne%20dual -o "out.json"

curl http://localhost:8080/autocomplete?term=Catheter%20balloon -o "out.json"

curl http://localhost:8080/autocomplete?term=Bandage%20elastoplast -o "out.json"

curl http://localhost:8080/autocomplete?term=virus -o "out.json"

curl http://localhost:8080/autocomplete?term=bacterial -o "out.json"

curl http://localhost:8080/autocomplete?term=ventricular -o "out.json"


## example output

C:\Data\github\nappi>go run .  
2025/01/21 19:10:33 Loading the NAPPI data...  
entry count: 454954  
Elapsed time: 340.2411ms  
2025/01/21 19:10:33 Building index for ProductName auto complete data...  
productNameIndex size: 288106  
Elapsed time: 201.7059ms  
2025/01/21 19:10:34 Server is running on port 8080...  
Elapsed time: 43.8917ms  
Elapsed time: 45.3861ms  
2025/01/21 19:11:05 Shutting down server...  
2025/01/21 19:11:05 Server exited cleanly

C:\Data\github\nappi>
