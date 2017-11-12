package frontend

const Html=`
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title> Cloud Conquerors </title>
    <link href="https://fonts.googleapis.com/css?family=Lato" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css?family=Oxygen" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css?family=Merriweather" rel="stylesheet">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.2/css/bootstrap.min.css" integrity="sha384-PsH8R72JQ3SOdhVi3uxftmaW6Vc51MKb0q5P2rRUpPvrszuE4W1povHYgTpBfshb" crossorigin="anonymous">
    <style>
		body {
    background-color: #DFE2DB;
    font-family: 'Oxygen', sans-serif;
}

.btn {
    font-family: 'Merriweather', serif;!important;
    margin-bottom: 1.5rem;
    transition: all 0.5s ease 0s;
    border-radius: 15px;
}

h1 {
    color: black;
    margin-top: 2rem;
    font-size: 3.5rem;
    margin-bottom: 2rem;
    text-shadow: 0 1px 0 #ccc,
    0 2px 0 #c9c9c9,
    0 3px 0 #bbb,
    0 4px 0 #b9b9b9,
    0 5px 0 #aaa,
    0 6px 1px rgba(0,0,0,.1),
    0 0 5px rgba(0,0,0,.1),
    0 1px 3px rgba(0,0,0,.3),
    0 3px 5px rgba(0,0,0,.2),
    0 5px 10px rgba(0,0,0,.25),
    0 10px 10px rgba(0,0,0,.2),
    0 20px 20px rgba(0,0,0,.15);
    font-family: 'Lato', sans-serif; !important;
}

.row {
    margin: 0.5rem 0 0.5rem 0;
}

	</style>
</head>

<body>
    <div class="container">
        <h1 class="text-center anim-typewriter"> Docker-Monitor </h1>
        <div class="row">
            <div class="col-sm-12 col-md-12 col-lg-2">
                <button type="button" class="btn btn-outline-primary"> Container Status </button>
            </div>
            <div class="col-sm-12 col-md-12 col-lg-8">
                <form>
                    <div class="form-group">
                        <div class="alert alert-success" role="alert">{{.Result}}</div>
                    </div>
                </form>
            </div>
        </div>

        <br>

        <div class="row">
            <div class="col-sm-12 col-md-12 col-lg-2">
                <button type="button" class="btn btn-outline-success"> Logs </button>
            </div>
            <div class="col-sm-12 col-md-12 col-lg-8">
                <form>
                    <div class="form-group">
                        <div class="alert alert-warning" role="alert">Updating soon ...</div>
                    </div>
                </form>
            </div>
        </div>
    </div>

    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.2/js/bootstrap.min.js" integrity="sha384-alpBpkh1PFOepccYVYDB4do5UnbKysX5WZXm3XxPqe5iKTfUKjNkCk9SaVuEZflJ" crossorigin="anonymous"></script>

</body>
</html>
`
