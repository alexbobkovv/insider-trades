import { Footer } from "components/Footer";
import { Header } from "components/Header";
import React from "react";

type Props = {
  children?: React.ReactNode;
};

export const Layout = ({ children }: Props) => {
  return (
    <>
      <Header />
      {children}
      <Footer />
    </>
  );
};
