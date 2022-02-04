import { Box, Button } from "@mui/material";
import { OptionsObject, SnackbarKey, SnackbarMessage } from "notistack";
import React, { FormEvent } from "react";
import ProgressBar from "./ProgressBar";

interface Props {
  handleFileUpload: (e: FormEvent<HTMLFormElement>) => void;
  setFile: React.Dispatch<React.SetStateAction<string | Blob | File>>;
  uploadProgress: number;
  uploading: boolean;
  file: string | Blob | File;
  enqueueSnackbar: (
    message: SnackbarMessage,
    options?: OptionsObject | undefined
  ) => SnackbarKey;
  closeSnackbar: (key?: SnackbarKey | undefined) => void;
}

const BoxView: React.FC<Props> = ({
  handleFileUpload,
  setFile,
  uploadProgress,
  uploading,
  file,
  enqueueSnackbar,
  closeSnackbar,
}) => {
  return (
    // box with margin on top
    <Box>
      <form onSubmit={handleFileUpload}>
        <h1>Select your database file</h1>
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
        <Button
          variant="contained"
          color="secondary"
          type="submit"
          disabled={!file || uploading}
        >
          Upload
        </Button>
        {uploadProgress > 0 && uploadProgress < 100 ? (
          <div className="progress-bar-container">
            <ProgressBar progress={uploadProgress} />
          </div>
        ) : null}
      </form>
    </Box>
  );
};

export default BoxView;
