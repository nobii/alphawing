package controllers

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/kayac/alphawing/app/models"
	"github.com/revel/revel"
)

type LimitedTimeController struct {
	AlphaWingController
	Bundle *models.Bundle
}

func (c *LimitedTimeController) GetDownloadPlist(bundleId int) revel.Result {
	bundle := c.Bundle

	app, err := bundle.App(c.Txn)
	if err != nil {
		panic(err)
	}

	ipaUrl, err := c.UriFor(fmt.Sprintf("bundle/%d/download_ipa", bundle.Id))
	if err != nil {
		panic(err)
	}

	p := models.NewPlist(app.Title, bundle.BundleVersion, ipaUrl.String())

	data, err := p.Marshall()
	if err != nil {
		panic(err)
	}

	c.Response.ContentType = "application/x-plist"
	return c.RenderText(string(data))
}

func (c *LimitedTimeController) GetDownloadIpa(bundleId int) revel.Result {
	resp, file, err := c.GoogleService.DownloadFile(c.Bundle.FileId)
	if err != nil {
		panic(err)
	}

	modtime, err := time.Parse(time.RFC3339, file.ModifiedDate)
	if err != nil {
		panic(err)
	}

	err = c.createAudit(models.ResourceBundle, bundleId, models.ActionDownload)
	if err != nil {
		panic(err)
	}

	c.Response.ContentType = "application/octet-stream"
	return c.RenderBinary(resp.Body, file.OriginalFilename, revel.Attachment, modtime)
}

func (c *LimitedTimeController) CheckNotFound() revel.Result {
	param := c.Params.Route["bundleId"]
	if 0 < len(param) {
		bundleId, err := strconv.Atoi(param[0])
		if err != nil {
			if err == sql.ErrNoRows {
				return c.NotFound("NotFound")
			}
			panic(err)
		}
		bundle, err := models.GetBundle(c.Txn, bundleId)
		if err != nil {
			panic(err)
		}
		c.Bundle = bundle
	}
	return nil
}

func (c *LimitedTimeController) CheckForbidden() revel.Result {
	if c.Bundle != nil {
		bundle := c.Bundle
		s, err := c.userGoogleService()
		if err != nil {
			panic(err)
		}
		_, err = s.GetFile(bundle.FileId)
		if err != nil {
			return c.Forbidden("Forbidden")
		}
	}
	return nil
}
