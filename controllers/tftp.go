package controllers

import (
	"errors"
	"fmt"
	"github.com/aarnaud/ipxeblue/utils"
	"github.com/pin/tftp/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetTFTPReader(config *utils.Config, db *gorm.DB) func(filename string, rf io.ReaderFrom) error {
	folder := "tftp"
	return func(filename string, rf io.ReaderFrom) error {
		filename = strings.TrimRight(filename, "�")
		raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()
		log.Info().Msgf("RRQ from %s filename %s", raddr.String(), filename)

		path := filepath.Join(folder, filename)

		// if file doesn't exist and path start by /grub/ use grubTFTP2HTTP
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) && strings.HasPrefix(filename, "/grub/") {
			return grubTFTP2HTTP(config, db, filename, rf)
		}

		file, err := os.Open(path)
		if err != nil {
			log.Error().Err(err)
			return err
		}
		stat, err := file.Stat()
		if err != nil {
			log.Error().Err(err)
			return err
		}
		rf.(tftp.OutgoingTransfer).SetSize(stat.Size())
		_, err = rf.ReadFrom(file)
		if err != nil {
			log.Error().Err(err)
			return err
		}
		return nil
	}
}

func GetTFTPWriter(config *utils.Config) func(filename string, wt io.WriterTo) error {
	return func(filename string, wt io.WriterTo) error {
		return nil
	}
}

func grubTFTP2HTTP(config *utils.Config, db *gorm.DB, filename string, rf io.ReaderFrom) error {
	gruburl, _ := url.Parse(config.BaseURL.String())
	if filename == "/grub/grub.cfg" {
		gruburl = gruburl.JoinPath("/grub/")
		resp, err := http.Get(gruburl.String())
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		reader := strings.NewReader(string(b))
		_, err = rf.ReadFrom(reader)
		if err != nil {
			return err
		}
		return nil
	}

	paths := strings.Split(filename, "/")
	if len(paths) < 11 {
		return fmt.Errorf("invalid path")
	}
	gruburl = gruburl.JoinPath("/grub/")
	query := gruburl.Query()
	query.Add("mac", paths[2])
	query.Add("ip", paths[3])
	query.Add("uuid", paths[4])
	query.Add("asset", strings.TrimSpace(strings.TrimLeft(paths[5], "-")))
	query.Add("manufacturer", strings.TrimLeft(paths[6], "-"))
	query.Add("serial", strings.TrimLeft(paths[7], "-"))
	query.Add("product", strings.TrimLeft(paths[8], "-"))
	query.Add("buildarch", strings.TrimLeft(paths[9], "-"))
	query.Add("platform", strings.TrimLeft(paths[10], "-"))
	gruburl.RawQuery = query.Encode()
	resp, err := http.Get(gruburl.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(b))
	_, err = rf.ReadFrom(reader)
	if err != nil {
		return err
	}
	return nil
}
