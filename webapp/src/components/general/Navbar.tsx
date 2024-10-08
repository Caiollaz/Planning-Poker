import { Button, Flex, Box, Heading, Spacer, useToast } from "@chakra-ui/react";
import Link from "next/link";
import { DarkModeSwitch } from "./DarkModeSwitch";
import { CopyIcon } from "@chakra-ui/icons";
import { useAppSelector } from "../../app/hooks";
import { BASE_WEBAPP_URL } from "../../constants";

export const Navbar = (props) => {
  const { isGame } = props;
  const { id } = useAppSelector((state) => state.room);

  if (isGame) {
    const toast = useToast();

    const handleCopyInvitation = () => {
      navigator.clipboard.writeText(`${BASE_WEBAPP_URL}/game/${id}`);
      return toast({
        position: "top",
        title: "Link de convite copiado para a área de transferência!",
        description: "Compartilhe com seus colegas de equipe e comece a jogar!",
        status: "success",
        variant: "subtle",
        duration: 3000,
        isClosable: true,
      });
    };

    return (
      <Flex position="fixed" p={5} w="100vw">
        <Flex w={{ lg: "95%", md: "95%", base: "85%" }}>
          <Box p="3">
            <Link href="/">
              <Heading size="md" _hover={{ cursor: "pointer" }}>
                Planning Poker
              </Heading>
            </Link>
          </Box>
          <Spacer />
          <Button
            leftIcon={<CopyIcon />}
            onClick={() => handleCopyInvitation()}
            mr="3rem"
            mt={1}
          >
            Copiar convite
          </Button>
        </Flex>
        <DarkModeSwitch />
      </Flex>
    );
  }

  return (
    <Flex position="fixed" p={5}>
      <Box p="3">
        <Link href="/">
          <Heading size="md" _hover={{ cursor: "pointer" }}>
            Planning Poker
          </Heading>
        </Link>
      </Box>
      <DarkModeSwitch />
    </Flex>
  );
};
