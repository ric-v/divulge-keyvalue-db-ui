import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";
import Brightness4Icon from "@mui/icons-material/Brightness4";
import Brightness7Icon from "@mui/icons-material/Brightness7";
import { Theme } from "@mui/material";

interface NavBarProps {
  toggleColorMode: () => void;
  theme: Theme;
}

const NavBar: React.FC<NavBarProps> = ({ toggleColorMode, theme }) => {
  return (
    <AppBar position="fixed">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Typography
            variant="h6"
            noWrap
            component="div"
            sx={{
              mr: 2,
              display: "flex",
              flexDirection: "row",
              justifyContent: "flex-start",
              alignItems: "end",
            }}
          >
            <img src="/logo-120px.png" width="40px" alt="logo" />
            <Typography variant="h4" ml={"1rem"}>
              Divulge
            </Typography>
            <Typography variant='subtitle2' ml={'12px'}> - Key Value DB Explorer</Typography>
          </Typography>

          <Box
            color={"body"}
            sx={{ flexGrow: 1, display: { xs: "none", md: "flex" } }}
          ></Box>

          <Box sx={{ display: "flex" }}>
            <IconButton
              sx={{ ml: 1 }}
              onClick={toggleColorMode}
              color="inherit"
            >
              {theme.palette.mode === "dark" ? (
                <Brightness7Icon />
              ) : (
                <Brightness4Icon />
              )}
            </IconButton>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};
export default NavBar;
