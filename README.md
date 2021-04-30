huginn - a REST backend [![Huginn](https://github.com/refundable-tgm/huginn/workflows/Huginn/badge.svg)](https://github.com/refundable-tgm/huginn/actions) [![codecov](https://codecov.io/gh/refundable-tgm/huginn/branch/master/graph/badge.svg?token=CKU7R0YUPC)](https://codecov.io/gh/refundable-tgm/huginn) [![Go Report Card](https://goreportcard.com/badge/github.com/refundable-tgm/huginn)](https://goreportcard.com/report/github.com/refundable-tgm/huginn)
=====

## Deployment

The deployment of this backend is managed by the install and infrastructure repository [refundable-tgm/install](https://github.com/refundable-tgm/install). This is the case because this backend doesn't get deployed alone, but with a frontend ([refundable-tgm/web](https://github.com/refudnable-tgm/web)) and the needed infrastructure to run Refundable.

## Documentation

After deployment a REST-API-Documentation provided through swagger will be accessible under `http://localhost:8080/`

## Components

This backend consists out of multiple components:
 
 - `assets`: containing pictures and other assets that are used during runtime
 - `db`: contains the API to connect and interact with the mongo database. It also includes the general data model.
 - `docs`: contains the swagger documentation of the REST-API created through [swaggo](https://github.com/swaggo/swag)
 - `excel_template`: contains the excel templates for the generation of travel invoices and business trip applications based on the official templates
 - `files`: contains the generation processes of all pdf and excel files and their pathings
 - `ldap`: contains tools to verify and get data from the TGM ldap service
 - `rest`: contains the actual REST-API with its endpoints, data structes, and token management
 - `untis`: contains the client to the WebUntis-API to interact with TGM's timetables

## Future Roadmap

 - [ ] implement file endpoints to just open existing files or just handle the pdf inside the application using byte slices
 - [ ] fix group lesson algorithm to group consecutive lessons  
 - [ ] implement sending of mails (when state changes or events occurr)
 - [ ] create logging system to log every event
 - [ ] implement the usage of existing applications as templates for new ones  
 - [ ] create simpler data model
 - [ ] change the design of the pdf templates to a more beautiful and easier to understand one
 - [ ] implement more control features as endpoints to be able to further interact with the data model instead of having to update it
 - [ ] create a moduler decision between databases including the API over a strategy pattern, to let the user decide which database to use while installing
 - [ ] create more checks for out of bounds values, to only allow valid values in the end
 - [ ] reduce cyclomatic complexity of code
 - [ ] introduce general performance improvements
 - [ ] design and implement further security measurements, like HTTPS for example

## Debug Mode

Debug mode of `gin-gonic` ([gin](https://github.com/gin-gonic/gin)) is automatically enabled when a `.debug` file is provided in `/vol/files/`

## Working Title

The working title under which this backend is developed is huginn. According to norse mythology Huginn and Muninn are the two ravens of Odin. Huginn translated into English means "to think", whereas Muninn means "to remember". As this backend symbolizes all "thinking" and processing done in this project this working title was chosen.


Grímnismál out of the Poetic Edda:
>O'er Mithgarth Huginn and Muninn both
>
>Each day set forth to fly;
>
>For Huginn I fear lest he come not home,
>
>But for Muninn my care is more.

Translation by Henry Adams Bellows