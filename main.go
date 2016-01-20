package main

import (
	// "fmt"
	"log"
	"os"	
	"github.com/codegangsta/cli"
	// "github.com/cheggaaa/pb"
	// "path/filepath"

)

func panicIfErr(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func main() {
	app := cli.NewApp()
  	app.Name = "Photo Manager"
  	app.Usage = "Organise your photo files"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		IndexCommand,
		DeduplicateCommand,
	}

	err := app.Run(os.Args)
	panicIfErr(err)

	// dirPath, e := os.Getwd()
	// check(e)

	// dbPath := filepath.Join(dirPath, ".db")

 // 	m, e := ReadMap(dbPath)

	// if e == nil {
	// 	fmt.Println("Index exists")
	// 	// check if the db file is up to date
	// 	dbTime, err := LastModTime(dbPath)
	// 	check(err)
	// 	dirTime, err := LastModTime(dirPath)
	// 	check(err)
		
	// 	if dbTime.Before(dirTime) {
	// 		fmt.Println("Index expired")
	// 		m = nil	
	// 	}		
	// } 	

	// if m == nil {		
	// 	fmt.Println("Counting files, this may take a few minutes...")
	// 	count, _ := CountFilesInPath(dirPath) 

	// 	fmt.Println("Indexing files...")
 
	// 	bar := pb.StartNew(count)

	// 	index := make(map[string][]string)
	// 	walkfunc := func(path string, fi os.FileInfo, err error) error {			
			
	// 		hash, err := ComputeMD5(path)

	// 		if err != nil {
	// 			return err
	// 		}

	// 		// check if index found 	
	// 		filePaths, ok := index[hash] 

	// 		if ok {				
	// 			index[hash] = append(filePaths, path)
	// 		} else {
	// 			index[hash] = []string{path}
	// 		}

	// 		bar.Increment()

	// 	 	return nil
	// 	}

	// 	err := WalkFilesInPath(dirPath, walkfunc)

	// 	check(err)

	// 	SaveMap(index, dbPath)
    	
	//     bar.FinishPrint("Index completed.")

	//     m = index
	// }

	// fmt.Println("Scanning duplications...")	

	// numDup := 0
	// for _, paths := range m {		
	// 	if len(paths) > 1 {			
	// 		numDup = numDup + len (paths) - 1			
	// 	}
	// }
	
	// if numDup > 0 {
	// 	fmt.Println("Removing duplicated files")	

	// 	bar := pb.StartNew(numDup)

	// 	for _, paths := range m {		
	// 		if len(paths) > 1 {			
	// 			// fileRemoved, err := file.RemoveAllButOne(paths)
	// 			// check(err)

	// 			// for i:=0;i<fileRemoved;i++ {
	// 				bar.Increment()
	// 			// }
	// 		}
	// 	}

	// 	bar.FinishPrint("Duplicated files removed.")

	// } else {
	// 	fmt.Println("No duplicated files found")			
	// }	 
	
	// // task 2: report duplications
	// // task 3: remove duplications
	// // optimise change detection. i.e. index last mod time for all dirs.
	// // maybe index all files, too
	// // task 4: move/rename files
	// // taks 5: clear empty directories
	// // task 6:

	//  fmt.Println(len(m))
}
