import Document, { Head, Html, Main, NextScript } from "next/document";

class MyDocument extends Document {
  render() {
    return (
      <Html lang="en">
        <Head>
          <link href="/favicon.ico" rel="shortcut icon" />
        </Head>
        <body className="main">
          <Main />
          <NextScript />
        </body>
      </Html>
    );
  }
}

export default MyDocument;
