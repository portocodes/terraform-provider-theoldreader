package main

import (
  "errors"
  "strings"
  "net/http"
  "encoding/json"
)

type SubscriptionList struct {
  Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
  Id string `json:"id"`
  Title string `json:"title"`
  URL string `json:"url"`
}

type QuickAdd struct {
  Query string `json:"query"`
  NumResults int `json:"numResults"`
  StreamId string `json:"streamId"`
  Error string `json:"error"`
}

type Client struct {
  Token string
}

func NewClient(token string) *Client {
  return &Client{ Token: token }
}

func (c *Client) Subscriptions() (SubscriptionList, error) {
  req, err := http.NewRequest("GET", "https://theoldreader.com/reader/api/0/subscription/list?output=json", nil)
  req.Header.Add("Authorization", "GoogleLogin auth=" + c.Token)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return SubscriptionList{}, err
  }

  var r SubscriptionList

  err = json.NewDecoder(resp.Body).Decode(&r)
  if err != nil {
    return SubscriptionList{}, err
  }

  return r, nil
}

func (c *Client) GetSubscription(id string) (Subscription, error) {
  r, err := c.Subscriptions()
  if err != nil {
    return Subscription{}, err
  }

  for _, subscription := range r.Subscriptions {
    if subscription.Id == id {
      return subscription, nil
    }
  }

  return Subscription{}, errors.New("unable to find subscription for id " + id)
}

func (c *Client) GetSubscriptionByURL(url string) (Subscription, error) {
  r, err := c.Subscriptions()
  if err != nil {
    return Subscription{}, err
  }

  for _, subscription := range r.Subscriptions {
    if subscription.URL == url {
      return subscription, nil
    }
  }

  return Subscription{}, errors.New("unable to find subscription for url " + url)
}

func (c *Client) CreateSubscription(url string) (Subscription, error) {
  req, err := http.NewRequest("POST", "https://theoldreader.com/reader/api/0/subscription/quickadd?quickadd=" + url, nil)
  req.Header.Add("Authorization", "GoogleLogin auth=" + c.Token)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return Subscription{}, err
  }

  var r QuickAdd

  err = json.NewDecoder(resp.Body).Decode(&r)
  if err != nil {
    return Subscription{}, err
  }

  if r.Error != "" {
    return Subscription{}, errors.New(r.Error)
  }

  return c.GetSubscriptionByURL(url)
}

func (c *Client) DeleteSubscription(id string) error {
  req, err := http.NewRequest("POST", "https://theoldreader.com/reader/api/0/subscription/edit", strings.NewReader("ac=unsubscribe&s="+id))
  req.Header.Add("Authorization", "GoogleLogin auth=" + c.Token)

  client := &http.Client{}
  _, err = client.Do(req)
  if err != nil {
    return err
  }

  return nil
}
