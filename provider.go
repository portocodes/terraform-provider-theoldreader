package main

import (
  "errors"
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/hashicorp/terraform/terraform"
  "github.com/hashicorp/terraform/plugin"
)

func resourceTheoldreaderSubscriptionCreate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*Client)
  url := d.Get("url").(string)

  subscription, err := client.CreateSubscription(url)
  if err != nil {
    return err
  }

  d.SetId(subscription.Id)
  return resourceTheoldreaderSubscriptionRead(d, meta)
}

func resourceTheoldreaderSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*Client)
  id := d.Id()

  subscription, err := client.GetSubscription(id)
  if err != nil {
    return err
  }

  d.Set("url", subscription.URL)
  d.Set("title", subscription.Title)
  return nil
}

func resourceTheoldreaderSubscriptionUpdate(d *schema.ResourceData, meta interface{}) error {
  return errors.New("oh no not implemented")
}

func resourceTheoldreaderSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*Client)
  id := d.Id()

  return client.DeleteSubscription(id)
}

func dataSourceTheoldreaderSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*Client)
  url := d.Get("url").(string)

  subscription , err := client.GetSubscriptionByURL(url)
  if err != nil {
    return err
  }

  d.SetId(subscription.Id)
  d.Set("url", subscription.URL)
  d.Set("title", subscription.Title)

  return nil
}

func Provider() terraform.ResourceProvider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "token": &schema.Schema{
        Type:        schema.TypeString,
        Description: "The Old Reader token",
        Required:    true,
        DefaultFunc: schema.EnvDefaultFunc("THE_OLD_READER_TOKEN", nil),
      },
    },
    ResourcesMap: map[string]*schema.Resource{
      "theoldreader_subscription": &schema.Resource{
        Schema: map[string]*schema.Schema{
          "url": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
          },
          "title": &schema.Schema{
            Type:     schema.TypeString,
            Computed: true,
          },
        },
        Create: resourceTheoldreaderSubscriptionCreate,
        Read:   resourceTheoldreaderSubscriptionRead,
        Update: resourceTheoldreaderSubscriptionUpdate,
        Delete: resourceTheoldreaderSubscriptionDelete,
        Importer: &schema.ResourceImporter{
          State: schema.ImportStatePassthrough,
        },
      },
    },
    DataSourcesMap: map[string]*schema.Resource{
      "theoldreader_subscription": &schema.Resource{
        Schema: map[string]*schema.Schema{
          "url": &schema.Schema{
            Type:     schema.TypeString,
            Required: true,
          },
          "title": &schema.Schema{
            Type:     schema.TypeString,
            Computed: true,
          },
        },
        Read: dataSourceTheoldreaderSubscriptionRead,
      },
    },
    ConfigureFunc: configureFunc(),
  }
}

func configureFunc() func(*schema.ResourceData) (interface{}, error) {
  return func(d *schema.ResourceData) (interface{}, error) {
    client := NewClient(d.Get("token").(string))
    return client, nil
  }
}

func main() {
  plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: func() terraform.ResourceProvider {
      return Provider()
    },
  })
}
