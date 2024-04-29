# Gocontentful

Gocontentful is a command line tool that generates a set of APIs for the [Go Language](https://go.dev) to interact with a [Contentful](https://www.contentful.com) CMS space.

Unlike the plain Contentful API for Go, the Gocontentful API is idiomatic. Go types are provided with names that mirror the content types of the Contentful space, and get/set methods are named after each field.

In addition, Gocontentful supports in-memory caching and updates of spaces. This way, the space is always accessible through fast Go function calls, even offline.

## What is Contentful

[Contentful](https://www.contentful.com/) is a content platform (often referred to as headless CMS) for [micro-content](https://www.contentful.com/r/knowledgebase/content-as-a-microservice/).

Unlike traditional CMSes, there's no pages or content trees in Contentful. The data model is built from scratch for the purpose of the consuming application, is completely flexible and can be created and hot-changed through the same Web UI that the content editors use. The model dictates which content types can reference others and the final structure is a graph.

## How applications interact with Contentful

Contentful hosts several APIs that remote applications use to create, retrieve, update and delete content. Content is any of the following:

- **Entries**, each with a content type name and a list of data fields as defined by the developer in the content model editor at Contentful
- **Assets** (images, videos, other binary files)

The Contentful APIs exist as either REST or GraphQL endpoints. Gocontentful only supports the REST APIs.

The REST APIs used to manage and retrieve content use standard HTTP verbs (GET, POST, PUT and DELETE) and a JSON payload for both the request (where needed) and the response.

## What is gocontentful

A golang API code generator that simplifies interacting with a Contentful space. The generated API:

- Supports most of the Contentful APIs to perform all read/write operation on entries and assets
- Hides the complexity of the Contentful REST/JSON APIs behind an idiomatic set of golang functions and methods
- Allows for in-memory caching of an entire Contentful space

## Why we need a Go API generator

While it's perfectly fine to call a REST service and receive data in JSON format, in Go that is not very practical. For each content type, the developer needs to maintan type definitions by hand and decode the JSON coming from the Contentful server into the value object.

In addition, calling a remote API across the Internet each time a piece of content is needed, even multiple times for a single page rendering, can have significant impact on performance.

Gocontentful generates a Go API that handles both issues above and can be regenerated every time the content model changes. The developer never needs to update the types by hand, or deal with the complexity of caching content locally. It all happens auytomatically in the generated client.

> **NOTE** - _How much code does Gocontentful generate? In a real-world production scenario where Gocontentful is in use as of 2024, a space content model with 43 content types of various field counts generates around 65,000 lines of Go code._
