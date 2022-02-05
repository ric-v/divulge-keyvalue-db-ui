import React, { FormEvent, useState } from "react";
import { Box, Button, Grid } from "@mui/material";
import { useSnackbar } from "notistack";
import http from "../services/axios-common";
import ProgressBar from "../components/ProgressBar";

type FileUploadProps = {
  setDbname: React.Dispatch<React.SetStateAction<string>>;
  setDbkey: React.Dispatch<React.SetStateAction<string>>;
  setStatus: React.Dispatch<React.SetStateAction<string>>;
};

const FileUpload = ({ setDbname, setDbkey, setStatus }: FileUploadProps) => {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [file, setFile] = useState<string | Blob | File>("");

  const [uploadProgress, updateUploadProgress] = useState(0);
  const [uploadStatus, setUploadStatus] = useState(false);
  const [uploading, setUploading] = useState(false);

  const handleFileUpload = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setUploading(true);
    const formData = new FormData();
    formData.append("file", file);

    http
      .post("/api/v1/upload?dbtype=buntdb", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        onUploadProgress: (ev: ProgressEvent) => {
          const progress = (ev.loaded / ev.total) * 100;
          updateUploadProgress(Math.round(progress));
        },
      })
      .then((resp) => {
        // our mocked response will always return true
        // in practice, you would want to use the actual response object
        setUploadStatus(true);
        setUploading(false);

        enqueueSnackbar("File uploaded successfully", {
          key: "success",
          variant: "success",
          draggable: true,
          onClick: () => {
            closeSnackbar("success");
          },
        });
        console.log(resp);
        setDbname(resp.data.filename);
        setDbkey(resp.data.dbkey);
        setStatus("connected");
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
      <Box
        component="form"
        onSubmit={handleFileUpload}
        sx={{
          "& > :not(style)": { m: 1 },
        }}
      >
        <Grid
          container
          columns={{ xs: 3, sm: 6, md: 8, lg: 12 }}
          rowSpacing={["1rem"]}
        >
          <Grid item xs={2} sm={4} md={6} lg={8}>
            <input
              type="file"
              name="file-choose"
              onChange={(e) => {
                e.preventDefault();
                if (e.target.files && e.target.files.length > 0) {
                  setFile(e.target.files[0]);
                  if (e.target.files[0].size > 1000000000) {
                    // 1GB
                    enqueueSnackbar("File size is too big", {
                      key: "file-size",
                      variant: "error",
                      onClick: () => {
                        closeSnackbar("file-size");
                      },
                    });
                    setFile("");
                  }
                }
              }}
            />
          </Grid>
          <Grid item xs={1} sm={2} md={2} lg={4}>
            <Button
              variant="contained"
              color="secondary"
              type="submit"
              disabled={!file || uploading}
            >
              Upload
            </Button>
          </Grid>
        </Grid>
        {uploadProgress > 0 && uploadProgress < 100 ? (
          <div className="progress-bar-container">
            <ProgressBar progress={uploadProgress} />
          </div>
        ) : null}
      </Box>
    </>
  );
};

export default FileUpload;
