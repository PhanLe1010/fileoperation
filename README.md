## A mini program to measure the time of file opening and closing operation

This program measure the time it take to open and close files. 
The result can be used to detect long file opening/closing operation.

### Installing and running
1. Install GoLang version >=1.20 https://go.dev/doc/install
2. `cd` into the correct directory that you want to test. For example, `cd /var/lib/longhorn`
3. Clone this repo `git clone https://github.com/PhanLe1010/fileoperation.git`
4. `cd fileoperation`
5. Run the program
    ```bash
    go run . [list-of-file-sizes] >> [output-file]
    ```
    for example,
    ```bash
    go run . 1 10 800 >> result.txt
    ```
    This command instructs the program to:
    1. Create and initiate (if not exist) 3 sparse files of size 1GB, 10GB,and 800GB. 
    Each file is filled with 512KB data chunks seperated by 512KB holes.
    So, the actual sizes are 0.5GB, 5GB, 400GB while the apparent sizes are 1GB, 10GB, 800GB
    2. Repeatedly opening and closing these files and record the time in the `result.txt` file

### Get result and cleanup
1. Let the program runs for a day or 2, then download and save the `result.txt` file
2. To clean up the test env, remove the git folder `sudo rm -r fileoperation`
