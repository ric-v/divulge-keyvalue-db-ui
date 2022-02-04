import React, { FormEvent, useState } from 'react';
import { Button } from "@mui/material";
import { useSnackbar } from 'notistack';
import http from '../services/axios-common';
import ProgressBar from './ProgressBar';

const Upload = () => {

  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [file, setFile] = useState<string | Blob | File>('');

  const [uploadProgress, updateUploadProgress] = useState(0);
  const [uploadStatus, setUploadStatus] = useState(false);
  const [uploading, setUploading] = useState(false);

  const handleFileUpload = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);

    http.post('/api/v1/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (ev: ProgressEvent) => {
        const progress = ev.loaded / ev.total * 100;
        updateUploadProgress(Math.round(progress));
      },
    })
      .then((resp) => {
        // our mocked response will always return true
        // in practice, you would want to use the actual response object
        setUploadStatus(true);
        setUploading(false);
        console.log(resp);

        enqueueSnackbar('success', {
          key: 'success',
          variant: 'success',
          draggable: true,
          onClick: () => {
            closeSnackbar('success');
          }
        });
      })
      .catch((err) => console.error(err));
  };

  return (
    <form onSubmit={handleFileUpload}>
      <h1>Select your files</h1>
      <input
        type='file'
        name='file-choose'
        onChange={(e) => {
          e.preventDefault();
          if (e.target.files && e.target.files.length > 0) {
            setFile(e.target.files[0])
            if (e.target.files[0].size > 1000000000) { // 1GB
              enqueueSnackbar('File size is too big', {
                key: 'file-size',
                variant: 'error',
                onClick: () => {
                  closeSnackbar('file-size');
                }
              });
              setFile('');
            }
          }
        }}
      />
      <Button
        variant='contained'
        color="secondary"
        type='submit'
        disabled={!file || uploading}
      >
        Upload
      </Button>
      {(uploadProgress > 0 && uploadProgress < 100)
        ?
        <div className="progress-bar-container">
          <ProgressBar progress={uploadProgress} />
        </div>
        : null
      }
    </form >
  );
};

export default Upload;
