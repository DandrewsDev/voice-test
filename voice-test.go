package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v4/pkg/demoinfocs/msgs2"
)

func extractOPUSData(inputVec [][]byte) ([]byte, error) {
	var combinedPLCData []byte

	for _, inputBytes := range inputVec {
		// Wrap the processing code in a try-except block
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Handle the error (print a message in this case)
					fmt.Println("Error occurred:", r)
				}
			}()

			// Check if there are enough bytes to read the payload type and length
			if len(inputBytes) < 3 {
				panic(fmt.Errorf("insufficient data length to read payload"))
			}

			// Read the payload type (byte 0)
			payloadType := inputBytes[0]

			// Read the length of OPUS PLC data as uint16 (bytes 1-2, assuming little-endian byte order)
			plcDataLength := binary.LittleEndian.Uint16(inputBytes[1:3])

			// Ensure there are enough bytes to read the entire OPUS PLC data
			if len(inputBytes) < int(plcDataLength)+3 {
				panic(fmt.Errorf("insufficient data length to read OPUS PLC data for ", inputBytes))
			}

			// Extract OPUS PLC data from inputBytes
			plcData := inputBytes[3 : 3+plcDataLength]

			// Print payload type, byte length, and OPUS PLC data
			fmt.Printf("[Payload Type: 0x%X] [Byte Length: %d] [OPUS PLC Data: %X]\n", payloadType, plcDataLength, plcData)

			// Append the extracted PLC data to the combined PLC data slice
			combinedPLCData = append(combinedPLCData, plcData...)
		}()
	}

	return combinedPLCData, nil
}

func createWAVFile(data []byte, filename string) error {
	fmt.Println("Saving wav file..")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := wav.NewEncoder(file, 48000, 16, 1, 1)

	buf := &audio.IntBuffer{
		Format: &audio.Format{
			SampleRate:  48000,
			NumChannels: 1,
		},
		Data: make([]int, len(data)/2), // Create a buffer for the audio data
	}

	// Convert the byte data to int samples
	for i := 0; i < len(buf.Data); i++ {
		buf.Data[i] = int(int16(binary.LittleEndian.Uint16(data[i*2 : (i+1)*2])))
	}

	if err := enc.Write(buf); err != nil {
		return err
	}
	return enc.Close()
}
func main() {
	var vec [][]byte

	f, err := os.Open("./1-48050850-ad9e-4497-bf8b-d20a7d7cb232.dem")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Configure parsing of BSPDecal net-message
	cfg := demoinfocs.DefaultParserConfig

	p := demoinfocs.NewParserWithConfig(f, cfg)
	defer p.Close()
	fmt.Println("Codec____________\n")
	p.RegisterNetMessageHandler(func(m *msgs2.CSVCMsg_VoiceInit) {
		fmt.Println(m.Codec, m, m.String(), m.Version)
	})
	fmt.Println("Data____________\n")
	p.RegisterNetMessageHandler(func(m *msgs2.CSVCMsg_VoiceData) {
		audioBytes := m.Audio.GetVoiceData()
		fmt.Println(hex.EncodeToString(audioBytes))
		parts := strings.Split(hex.EncodeToString(audioBytes), "0bc05d")

		// Remove the first part and append the rest to vec

		if len(parts) > 1 {
			remainingHex := parts[1]
			decodedBytes, err := hex.DecodeString(remainingHex)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			vec = append(vec, decodedBytes)
		}

		fmt.Println("========================")
	})
	// Print the resulting vec

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}

	// Extract OPUS PLC data
	opusData, err := extractOPUSData(vec)
	if err != nil {
		fmt.Println("Invalid input data:", err)
		return
	}

	// Create WAV file from OPUS data
	err = createWAVFile(opusData, "output.wav")
	if err != nil {
		log.Fatalf("Error creating WAV file: %v", err)
		return
	}

	fmt.Println("Audio file generated: output.wav")

}
