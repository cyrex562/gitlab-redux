import Vue from 'vue';
import store from '~/mr_notes/stores';
import { pinia } from '~/pinia/instance';
import { DiffFile } from '~/rapid_diffs/diff_file';
import FileBrowser from './file_browser.vue';

export function initFileBrowser() {
  const el = document.querySelector('[data-file-browser]');
  // eslint-disable-next-line no-new
  new Vue({
    el,
    store,
    pinia,
    render(h) {
      return h(FileBrowser, {
        on: {
          clickFile(file) {
            DiffFile.findByFileHash(file.fileHash).selectFile();
          },
        },
      });
    },
  });
}
