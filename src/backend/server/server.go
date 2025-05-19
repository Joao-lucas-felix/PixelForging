package server

import (
	"io"
	"log"
	"net"

	pixelforging_grpc "github.com/Joao-lucas-felix/PixelForging/src/backend/pb/pixelforging-grpc"
	pixelforging "github.com/Joao-lucas-felix/PixelForging/src/image-processing"
	"google.golang.org/grpc"
)

type Server struct {
	pixelforging_grpc.PixelForgingServer
}

func (s Server) ExtractPalette(srv pixelforging_grpc.PixelForging_ExtractPaletteServer) error {
	var pixelArt []byte
	var fileName, fileType string
	var colorsPerRow, colorWidth, colorHeight, colorNum int32

	log.Println("Extracting palette...")
	for {
		data, err := srv.Recv()
		if err == io.EOF {
			log.Println("Finished receiving data")
			break
		}
		if err != nil {
			log.Println("Error receiving data: ", err)
			return err
		}

		pixelArt = append(pixelArt, data.GetFileBytes()...)
		fileName = data.GetFileName()
		fileType = data.GetFileType()

		// Get the color palette parameters
		colorsPerRow = data.GetColorsPerRow()
		colorHeight = data.GetColorHeight()
		colorWidth = data.GetColorWidth()
		colorNum = data.GetColorNum()
	}
	log.Println("Received file:\t", fileName)
	log.Println("File type:\t", fileType)
	log.Println("File size:\t", len(pixelArt))

	// Convert bytes to image
	img, _, err := pixelforging.BytesToImage(pixelArt, fileType)
	if err != nil {
		log.Println("Error converting bytes to image: ", err)
		return err
	}
	// Extract palette from image
	img = pixelforging.ExtractColorPalette(img, int(colorsPerRow), int(colorWidth), int(colorHeight), int(colorNum))
	
	bytesOutput, err := pixelforging.ImageToBytes(img, fileType)
	if err != nil {
		log.Println("Error converting image to bytes: ", err)
		return err
	}
	// Implementar a logica de chunks
	log.Println("Palette extracted successfully")
	log.Println("Sending data...")
	for _, bytes := range bytesOutput {
		err := srv.Send(&pixelforging_grpc.ExtractPaletteOutput{
			PaletteBytes: []byte{bytes},
			FileName:     fileName,
			FileType:     fileType,
		})
		if err != nil {
			log.Println("Error sending data: ", err)
			return err
		}
	}
	return nil
}

// BoostrapServer starts the gRPC server on port 9090
// @Description: Starts the gRPC server on port 9090
func BoostrapServer(port string) {
	log.Println("Starting gRPC server...")
	// Create a new gRPC server

	listner, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln("Error starting server: ", err)
	}
	grpcServer := grpc.NewServer()
	pixelforging_grpc.RegisterPixelForgingServer(grpcServer, &Server{})

	log.Println("Server started on port 9090")
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalln("Error starting server: ", err)
	}

}
