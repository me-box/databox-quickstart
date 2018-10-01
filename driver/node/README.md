##Writing a driver in node

This driver will write a value to the key value store and then access it when the page is refreshed.  To get started, first ensure that your are running the test environment (i.e. you have called [home]/databox-quickstart/testenv/start.sh). Then run:

```
cd src
npm install
npm run testmode
```

Then go to http://127.0.0.1:8080.  In the input box type some text and hit update.  Now refresh the page and you should see a statement "current config is: [your text]".  This means the driver has successfully set up a store to read/write to. 