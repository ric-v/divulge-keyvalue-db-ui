import { useState, useEffect } from "react";
import {
  Tooltip,
  Fab,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  TextField,
  DialogActions,
  Button,
  IconButton,
  List,
  ListItemText,
  ListItemButton,
} from "@mui/material";
import http from "../services/axios-common";
import AddIcon from "@material-ui/icons/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import { useSnackbar } from "notistack";
import { Box } from "@mui/system";

type Props = {
  defBucket: string;
  setBucket: React.Dispatch<React.SetStateAction<string>>;
};

const ManageBucketsModal = ({ defBucket, setBucket }: Props) => {
  const [open, setOpen] = useState(false);
  const [addBucket, setAddBucket] = useState("");
  const [buckets, setBuckets] = useState<string[]>([]);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  // use effect to get buckets
  useEffect(() => {
    loadBuckets();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // load buckets from server
  const loadBuckets = () => {
    http(String(localStorage.getItem("dbkey")))
      .get("/api/v1/bucket")
      .then((resp) => {
        setBuckets(resp.data.data.buckets);
        setBucket(resp.data.data.defaultBucket);
        setAddBucket("");
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

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setAddBucket("");
  };

  // add bucket
  const handleSubmit = () => {
    console.log("bucket name: ", addBucket);

    http(String(localStorage.getItem("dbkey")))
      .post("/api/v1/bucket?bucket=" + addBucket)
      .then((resp) => {
        enqueueSnackbar("Bucket Added Successfully!", {
          key: addBucket,
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar(addBucket);
          },
        });
        console.log(resp);
        setAddBucket("");
        loadBuckets();
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

  const handleSetBucket = (bucket: string) => {
    http(String(localStorage.getItem("dbkey")))
      .put("/api/v1/bucket/" + bucket)
      .then((resp) => {
        enqueueSnackbar("Default Bucket Selected!", {
          key: bucket,
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar(bucket);
          },
        });
        console.log(resp);
        setAddBucket("");
        loadBuckets();
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

  const handleDelete = (bucket: string) => {
    http(String(localStorage.getItem("dbkey")))
      .delete("/api/v1/bucket?bucket=" + bucket)
      .then((resp) => {
        enqueueSnackbar("Bucket Removed Successfully!", {
          key: bucket,
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar(bucket);
          },
        });
        console.log(resp);
        setAddBucket("");
        loadBuckets();
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
    <div style={{ marginLeft: "20px" }}>
      <Tooltip title="Add / Update / Remove Buckets" placement="top">
        <Fab
          variant="extended"
          color="secondary"
          onClick={() => handleClickOpen()}
        >
          {/* <AddIcon /> */}
          <Typography ml={1}>Manage Buckets</Typography>
        </Fab>
      </Tooltip>
      <Dialog
        open={open}
        onClose={handleClose}
        sx={{
          backdropFilter: "blur(3px)",
        }}
      >
        <DialogTitle>Manage Bolt DB Buckets</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Add / Remove / Update and set Default bucket for the Bolt DB.
          </DialogContentText>
          <Box mt={2}>
            <Typography sx={{ mt: 4, mb: 2 }} variant="h6" component="div">
              Set default bucket
            </Typography>
            {/* <Demo> */}
            <List
              sx={{
                bgcolor: "background.paper",
              }}
            >
              {buckets ? (
                buckets?.map((bucket) => (
                  <ListItemButton
                    role={undefined}
                    onClick={() => handleSetBucket(bucket)}
                    dense
                    sx={{
                      "&:hover": {
                        backgroundColor: "primary.main",
                        color: "primary.contrastText",
                      },
                      backgroundColor:
                        bucket === defBucket ? "secondary.main" : "",
                      color:
                        bucket === defBucket ? "secondary.contrastText" : "",
                    }}
                  >
                    <ListItemText
                      id={bucket}
                      primary={bucket}
                      secondary={defBucket === bucket ? "(Default)" : ""}
                    />
                    <IconButton
                      edge="end"
                      aria-label="delete"
                      onClick={() => handleDelete(bucket)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </ListItemButton>
                ))
              ) : (
                <Typography color={"error"}>
                  No Buckets. Add buckets to continue.
                </Typography>
              )}
            </List>
            {/* </Demo> */}
          </Box>
          <Box marginTop={2}>
            <TextField
              autoFocus
              margin="dense"
              id="add-bucket"
              label="Add New Bucket"
              value={addBucket}
              type="text"
              fullWidth
              variant="outlined"
              onChange={(e) => setAddBucket(e.target.value)}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleSubmit}>Add</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default ManageBucketsModal;
