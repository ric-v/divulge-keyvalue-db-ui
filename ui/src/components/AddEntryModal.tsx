import { useState } from "react";
import AddIcon from "@material-ui/icons/Add";
import {
  Tooltip,
  Fab,
  Typography,
  Button,
  TextField,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";
import http from "../services/axios-common";
import { useSnackbar } from "notistack";

type AddEntryModalProps = {
  setUpdated: React.Dispatch<React.SetStateAction<boolean>>;
};

const AddEntryModal = ({ setUpdated }: AddEntryModalProps) => {
  const [open, setOpen] = useState(false);
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setKey("");
    setValue("");
  };

  const handleSubmit = () => {
    console.log(key, value);

    const payload = {
      key,
      value,
    };

    http(String(localStorage.getItem("dbkey")))
      .post("/api/v1/db", payload)
      .then((resp) => {
        enqueueSnackbar("File created successfully", {
          key: "success",
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar("success");
          },
        });
        console.log(resp);
        setUpdated(true);
        handleClose();
      })
      .catch((err) => {
        enqueueSnackbar(err.response.data.message, {
          key: "error",
          variant: "error",
          draggable: true,
          onClick: () => {
            closeSnackbar("error");
          },
        });
      });
  };

  return (
    <>
      <Tooltip title="Add New Entry" placement="top">
        <Fab
          variant="extended"
          color="primary"
          onClick={() => handleClickOpen()}
        >
          <AddIcon />
          <Typography ml={1}>Add</Typography>
        </Fab>
      </Tooltip>
      <Dialog
        open={open}
        onClose={handleClose}
        sx={{
          backdropFilter: "blur(3px)",
        }}
      >
        <DialogTitle>Add New Key-Value Pair</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Insert new key-value pair to the DB. Key must be unique if not the
            existing record will be overwritten, and the value will be updated.
            Value can be any text type data.
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="key"
            label="Key"
            value={key}
            type="key"
            fullWidth
            variant="standard"
            onChange={(e) => setKey(e.target.value)}
          />
          <TextField
            margin="dense"
            id="value"
            label="value"
            minRows={5}
            multiline
            value={value}
            type="value"
            fullWidth
            variant="standard"
            onChange={(e) => setValue(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleSubmit}>Add</Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default AddEntryModal;
