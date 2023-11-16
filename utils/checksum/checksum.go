package checksum

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

type Checksum struct {
	Input     string
	InputList string
}

// Generate calculates the MD5 checksum of the file located at filePath and returns it as a hexadecimal string.
// If an error occurs while opening or reading the file, it returns an empty string and the error.
func Generate(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create an MD5 hash object
	hasher := md5.New()

	// Copy the file content to the hash object
	_, _ = io.Copy(hasher, file)

	// Get the checksum as a byte slice
	checksumBytes := hasher.Sum(nil)

	// Convert the checksum to a hexadecimal string
	checksum := hex.EncodeToString(checksumBytes)

	return checksum, nil
}

// Read reads the file at the given file path and returns a Checksum struct
// containing the first two lines of the file as Input and InputList fields.
func Read(filePath string) (Checksum, error) {

	c := Checksum{}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return c, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the first two lines
	for i := 0; i < 2 && scanner.Scan(); i++ {
		if i == 0 {
			c.Input = scanner.Text()
		} else {
			c.InputList = scanner.Text()
		}
	}

	// Check for errors during scanning
	// if err := scanner.Err(); err != nil {
	// 	return Checksum{}, err
	// }

	return c, nil
}

// AreEqual checks if two Checksum structs are equal by comparing their Input and InputList fields.
func AreEqual(checksum1, checksum2 Checksum) bool {
	return checksum1.Input == checksum2.Input && checksum1.InputList == checksum2.InputList
}

// Write writes the given Checksum to the file at the given filePath.
// It returns an error if the file cannot be created or if there is an error while writing to the file.
func Write(filePath string, c Checksum) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the checksums
	_, _ = file.WriteString(c.Input + "\n" + c.InputList)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// ValidateChecksum validates the checksum of the input files and returns a boolean indicating if the checksum is valid or not.
// If there is a checksum.txt file in the workdir, then this is a restart run.
// If there is no checksum.txt file in the workdir, then this is a fresh run.
// If there is a checksum.txt file in the workdir, but the checksum is different, an error is returned.
// The function takes three parameters:
// - inputF: the path to the input file.
// - inputL: the path to the input list file.
// - workdir: the path to the working directory.
// The function returns a boolean indicating if the checksum is valid or not, and an error if any.
func ValidateChecksum(inputF, inputL, checksumF string) (bool, error) {

	c := Checksum{}
	v, errGenChecksumInput := Generate(inputF)
	if errGenChecksumInput != nil {
		return false, errGenChecksumInput
	}
	c.Input = v

	v, errGenChecksumInputList := Generate(inputL)
	if errGenChecksumInputList != nil {
		return false, errGenChecksumInputList
	}
	c.InputList = v

	// Check if the checksum file exists
	_, err := os.Stat(checksumF)
	if os.IsNotExist(err) {
		errWriteChecksum := Write(checksumF, c)
		if errWriteChecksum != nil {
			return false, errWriteChecksum
		}
		return false, nil
	}

	// Load the checksum file
	oldChecksum, _ := Read(checksumF)
	// if err != nil {
	// 	return false, err
	// }

	if !AreEqual(c, oldChecksum) {
		return false, errors.New("the input files have changed since the last run. Remove " + checksumF + " to force a fresh run.")
	} else {
		return true, nil
	}
}
