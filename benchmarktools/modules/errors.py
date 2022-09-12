class Error(Exception):
    """Base class for other exceptions"""

    pass


class ConfigKeyUndefinedError(Error):
    """Raised when a mandatory key is not defined in the config file."""

    def __init__(self, config_key, message="Config key not defined"):
        self.config_key = config_key
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.config_key}"


class ConfigKeyEmptyError(Error):
    """Raised when a key is empty in the config file."""

    def __init__(self, empty_config_key, message="Config key is empty"):
        self.empty_config_key = empty_config_key
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.empty_config_key}"


class PathNotFound(Error):
    """Raised when a Path is not found."""

    def __init__(self, path, message="Path not found"):
        self.path = str(path)
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.path}"


class HeaderUndefined(Error):
    """Raised when a necessary Header is not defined in the config file."""

    def __init__(self, header_name, message="Header not found"):
        self.header_name = header_name
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.header_name}"


class ScenarioUndefined(Error):
    """Raised when no scenarios have been defined in the config file."""

    def __init__(self, message="No scenarios have been found"):
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}"


class InvalidParameter(Error):
    """Raised when an invalid parameter is found."""

    def __init__(self, parameter, message="Parameter invalid"):
        self.parameter = parameter
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.parameter}"


class InvalidRunName(Error):
    """Raised when a run name is not valid."""

    def __init__(self, runname, message="Run name invalid"):
        self.runname = runname
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.runname}"


class SuffixError(Error):
    """Raised when the input files do not match the suffix."""

    def __init__(self, target, suffix, message="No PDB input matches the suffix"):
        self.suffix = suffix
        self.target = target
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.suffix} at folder {self.target}/"


class MultipleInputError(Error):
    """Raised when multiple input match the same suffix."""

    def __init__(
        self, target, suffix, message="More than one PDB input matches the suffix"
    ):
        self.suffix = suffix
        self.target = target
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.suffix} at {self.target}"


class HaddockError(Error):
    """Raised when HADDOCK can not be executed."""

    def __init__(self, output_file, message="HADDOCK could not be executed"):
        self.output_file = output_file
        self.message = message
        super().__init__(self.message)

    def __str__(self):
        return f"{self.message}: {self.output_file}"
