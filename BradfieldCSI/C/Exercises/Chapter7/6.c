/* Write a program to compare two files, print the first line where they differ*/


/* Plan:
1. Define a maxline length to be read
2. Accept/Expect 2 file names via command line arguments (handle errors if > 2 or < 2)
3. Open both using fopen in read mode
4. Read 1 line each for both - then compare - if match then continue, if no match/one file over handle
5. fclose - handled by exit in all cases
*/

/*Note: 
- If number of lines differs in both files, even if the text is same until one of the file terminates -> will mark as different as one of the files will terminate with a \n while the other will not
- Can handle with a custom comparer where we first compute length of the result of fgets each time and then handle last character accordingly
*/

#include<stdio.h>
#include<stdlib.h>
#include<string.h>

#define MAXLINELENGTH 20

int main(int argc, char *argv[]) {

    FILE *fp1, *fp2;
    char linef1[MAXLINELENGTH], linef2[MAXLINELENGTH];
    char *fgets1Res, *fgets2Res; // to store results for fgets, for error handling later
    

    if (argc != 3) { // > 2 or < 2 file names
        fprintf(stderr, "%d arguments provided as filenames when 2 were expected", argc-1);
        exit(1);
    }

    fp1 = fopen(*++argv, "r");
    fp2 = fopen(*++argv, "r");

    if (fp1 == NULL || fp2 == NULL) { // file opening error
        fprintf(stderr, "Error opening file with name provided in argument number %d", fp1 == NULL ? 1 : 2);
        exit(2);
    }

    while ((fgets1Res = fgets(linef1,MAXLINELENGTH, fp1)) && (fgets2Res = fgets(linef2, MAXLINELENGTH, fp2))) { // comparing line by line
        if (strcmp(fgets1Res, fgets2Res) != 0) {
            fprintf(stdout, "The first line(s) of difference between the two files are:\n - %s \n - %s \n respectively.", linef1, linef2);
            exit(0);
        }
    }

    if (ferror(fp1) || ferror(fp2)) { // error in file stream
        fprintf(stderr, "Stream error in file with name provided in argument number %d", ferror(fp1) ? 1 : 2);
        exit(3);
    }

    // if ((fgets1Res && !fgets2Res) || (!fgets1Res && fgets2Res)) { // cannot handle case where one file ends and the other does not - read 'Note'
    //     fprintf(stdout, "One file ended and the other did not - last line read in file provided by arg %d is %s", (fgets1Res && !fgets2Res) ? 1 : 2, (fgets1Res && !fgets2Res) ? linef1 : linef2);
    //     exit(0);
    // }

    fprintf(stdout, "No line of difference found between the two files");
    exit(0);





    



    return 0;
}