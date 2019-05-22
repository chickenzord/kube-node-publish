package dns

import (
	"github.com/cloudflare/cloudflare-go"
)

func EnsureRecord(domainName string, recordType string, recordContent string) error {
	api, err := cloudflare.New(cfAPIKey, cfEmail)
	if err != nil {
		return err
	}

	zoneID, err := api.ZoneIDByName(cfZone)
	if err != nil {
		return err
	}

	newRecord := cloudflare.DNSRecord{
		Name:    domainName,
		Type:    "A",
		Content: recordContent,
	}

	records, err := api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: domainName, Type: "A"})
	if err != nil {
		return err
	}

	if len(records) > 0 {
		oldRecord := records[0]
		if oldRecord.Content != newRecord.Content {
			if err := api.UpdateDNSRecord(zoneID, oldRecord.ID, newRecord); err != nil {
				return err
			}
		}
	} else {
		resp, err := api.CreateDNSRecord(zoneID, newRecord)
		if err != nil || !resp.Success {
			return err
		}
	}

	return nil
}
