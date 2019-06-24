# Grafana dashboard

## Importing the dashboard

The *Kubernetes CodeNotary vcn Overview* dashboard can be imported into Grafana by copy-paste the provided [dashboard.json](dashboard.json) or by [getting it from the Grafana Dashboards](https://grafana.com/dashboards/10339) website.

You can find detailed Grafana importing instructions [here](https://grafana.com/docs/reference/export_import/). 


## Creating alerts

Once the provided Grafana dashboard has been imported, you can edit the *Status Count* panel:

![Step 1: Edit the Status Count panel](img/alert-step1-edit.png?raw=true "Edit the Status Count panel")

Then select the *Alert* tab and click on the *Create Alert* button:

![Step 2: Create the alert](img/alert-step2-create-alert.png?raw=true "Create the alert")

Finally, customize your alert as shown in the following example:

![Step 3: Customize the alert](img/alert-step3-example.png?raw=true "Customize the alert")

The *Status Count* pannel comes with four queries that you can use within alert conditions:

- Query **A** count of *Trusted* images
- Query **B** count of *Unknown* images
- Query **C** count of *Unsupported* images
- Query **D** count of *Untrusted* images