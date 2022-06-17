import { Footer } from "components/Footer";
import { Header } from "components/Header";
import Navbar from "components/Navbar";
import SearchBar from "components/SearchBar";
import React from "react";

type Props = {
  children?: React.ReactNode;
  // title?: string
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
