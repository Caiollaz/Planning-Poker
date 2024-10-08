import { Box, Button } from "@chakra-ui/react";
import { Navbar } from "../components/general/Navbar";
import { Container } from "../components/general/Container";
import { Hero } from "../components/general/Hero";
import Link from "next/link";

function Kicked() {
  return (
    <Box height="100vh">
      <Navbar />
      <Container height="100%">
        <Hero
          title="Você foi expulso! 😵"
          subText="Parece que o administrador não gosta mais de você."
        >
          <Link href="/">
            <Button variant="solid" size="lg">
              Voltar para o início
            </Button>
          </Link>
        </Hero>
      </Container>
    </Box>
  );
}

export default Kicked;
