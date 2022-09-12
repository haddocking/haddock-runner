import copy
import logging
import os
import re
import shutil
import sys
from pathlib import Path

from benchmarktools.modules.errors import HaddockError, SuffixError

datasetlog = logging.getLogger("setuplog")


def glob_pdb_re(target_path, regex):
    """Apply a regex in Path.glob function"""
    found_list = []
    for pdb_path in target_path.glob("*pdb"):
        matched = re.findall(regex, str(pdb_path))
        if matched:
            found_list.append(pdb_path)
    return found_list


class Dataset:
    def __init__(self, dataset_path):
        self.dataset_path = dataset_path

    def check_input_files(self, receptor_suffix, ligand_suffix):
        """Check if the files inside the dataset are formatted correctly."""
        # glob returns a abritrary ordered list
        target_directory_list = self.dataset_path.glob("*")

        for target in target_directory_list:
            if os.path.isfile(target):
                # this is a file, dont check it
                continue

            # Demo of this regex using receptor_suffix = r_u
            # https://regex101.com/r/8MviAN/1
            receptor_regex = rf"(.*{receptor_suffix}_?)(\d.pdb|.pdb)"
            receptor_l = glob_pdb_re(target, receptor_regex)

            # FIXME: for CG mode to work we need to rewrite this part
            receptor_l = self.remove_cg(receptor_l)

            if len(receptor_l) == 0:
                raise SuffixError(target.name, receptor_suffix)
            elif not receptor_l:
                raise SuffixError(target.name, receptor_suffix)

            # Demo of this regex using ligand_suffix = l_u
            # https://regex101.com/r/8n3Bg3/1
            ligand_regex = rf"(.*{ligand_suffix}_?)(\d.pdb|.pdb)"
            ligand_l = glob_pdb_re(target, ligand_regex)

            # FIXME: for CG mode to work we need to rewrite this part
            ligand_l = self.remove_cg(ligand_l)

            if len(ligand_l) == 0:
                raise SuffixError(target.name, ligand_suffix)
            elif not ligand_l:
                raise SuffixError(target.name, ligand_suffix)

    def setup(
        self,
        haddock,
        parameters,
        receptor_suffix,
        ligand_suffix,
        setup_mode=False,
        force=False,
    ):
        """Setup the HADDOCK run."""

        setup_paths = []
        for target in self.dataset_path.glob("*"):
            if os.path.isfile(target):
                # its a file again
                continue

            datasetlog.debug(f"Setting up {target}")

            setup_run_name = target / f'run-{parameters["run_name"]}'

            if setup_run_name.exists():
                datasetlog.warning(f"Run folder already exist: {setup_run_name}")
                if force:
                    datasetlog.warning(f"Force removing {setup_run_name}")
                    shutil.rmtree(setup_run_name)
                else:
                    datasetlog.error(
                        "Resuming is not yet supported, please "
                        "remove all the run directories and "
                        "restart"
                    )
                    sys.exit()

            # Demo of this regex using receptor_suffix = r_u
            # https://regex101.com/r/8MviAN/1
            receptor_regex = rf"(.*{receptor_suffix}_?)(\d.pdb|.pdb)"
            receptor_l = glob_pdb_re(target, receptor_regex)
            receptor_l = self.remove_cg(receptor_l)
            # Warning: This sorting is important!
            receptor_l.sort()

            # Demo of this regex using ligand_suffix = l_u
            # https://regex101.com/r/8n3Bg3/1
            ligand_regex = rf"(.*{ligand_suffix}_?)(\d.pdb|.pdb)"
            ligand_l = glob_pdb_re(target, ligand_regex)
            ligand_l = self.remove_cg(ligand_l)
            # Warning: This sorting is important!
            ligand_l.sort()

            param_str = ""
            if "ambig_tbl" in parameters:
                param_str += f'AMBIG_TBL={parameters["ambig_tbl"]}' + os.linesep

            param_str += f"PDB_FILE1={receptor_l[0]}" + os.linesep
            param_str += f"PDB_FILE2={ligand_l[0]}" + os.linesep

            if len(receptor_l) > 1:
                # This is an ensemble, generate a list
                list_name = (
                    receptor_l[0].name.split(receptor_suffix)[0]
                    + receptor_suffix
                    + ".list"
                )
                list_file = Path(target, list_name)
                with open(list_file, "w") as fh:
                    for receptor in receptor_l:
                        # this needs to be the absolutepath
                        fh.write(f'"{receptor}"{os.linesep}')
                fh.close()
                param_str += f"PDB_LIST1={list_file}" + os.linesep

            if len(ligand_l) > 1:
                # This is an ensemble, generate a list
                list_name = (
                    ligand_l[0].name.split(ligand_suffix)[0] + ligand_suffix + ".list"
                )
                list_file = Path(target, list_name)
                with open(list_file, "w") as fh:
                    for ligand in ligand_l:
                        # this needs to be the absolutepath
                        fh.write(f'"{ligand}"{os.linesep}')
                fh.close()
                param_str += f"PDB_LIST2={list_file}" + os.linesep

            param_str += f"PROJECT_DIR={target}" + os.linesep
            param_str += "N_COMP=2" + os.linesep
            param_str += f'RUN_NUMBER=-{parameters["run_name"]}' + os.linesep
            param_str += f"HADDOCK_DIR={haddock.path}" + os.linesep

            # TODO: Get the segID from the configuration file
            param_str += "PROT_SEGID_1=A" + os.linesep
            param_str += "PROT_SEGID_2=B" + os.linesep

            run_param_fname = target / f'run.param-{parameters["run_name"]}'

            with open(run_param_fname, "w") as param_fh:
                param_fh.write(param_str)

            shutil.copy(run_param_fname, run_param_fname.parent / "run.param")

            try:
                lig_top = list(target.glob("ligand.top"))[0]
            except IndexError:
                lig_top = False

            try:
                lig_par = list(target.glob("ligand.param"))[0]
            except IndexError:
                lig_par = False

            try:
                haddock.setup(target_dir=target, identifier=parameters["run_name"])
            except HaddockError as e:
                datasetlog.error(e)
                sys.exit()

            # change parameters in run.cns
            run_cns = setup_run_name / "run.cns"

            # save the original run.cns
            run_cns_ori = run_cns.parent / "run.cns-ori"
            shutil.copy(run_cns, run_cns_ori)

            # copy the edited
            edited_cns = self.edit_cns(run_cns, parameters)
            shutil.copy(edited_cns, run_cns)

            # copy the ligands inside the run, if any
            if lig_top and lig_par:
                datasetlog.debug("Ligand param/top found, adding it to run")
                shutil.copy(lig_top, setup_run_name / "toppar/ligand.top")
                shutil.copy(lig_par, setup_run_name / "toppar/ligand.param")

            # Setup done, add its path to a list
            setup_paths.append(setup_run_name)

        return setup_paths

    @staticmethod
    def edit_cns(run_cns, parameters):
        """Edit the run.cns parameter file."""
        param_regex = r"{===>}\s(\w*)=(.*)\;"
        edited_cns = run_cns.parent / "run.cns-edit"

        with open(edited_cns, "w") as edited_cns_fh:
            with open(run_cns, "r") as cns_fh:

                for line in cns_fh.readlines():
                    if line.startswith("{===>}"):

                        param, value = re.findall(param_regex, line)[0]

                        if param in parameters:
                            custom_value = parameters[param]

                            # workaround to handle booleans
                            custom_value = str(custom_value).lower()

                            if custom_value != value:
                                datasetlog.debug(
                                    f"Changing {param} from {value} to {custom_value}"
                                )

                                line = f"{{===>}} {param}={custom_value};" + os.linesep

                    edited_cns_fh.write(line)

        edited_cns_fh.close()

        return edited_cns

    @staticmethod
    def remove_cg(path_list):
        """Removes _cg.pdb from a list of Posix paths."""
        clean_list = copy.deepcopy(path_list)
        for path in clean_list:
            if "cg" in str(path):
                clean_list.remove(path)
        return clean_list
