
package domain

type Image struct {
	id string
	fileName string
	contentFile string
	owner string
	size float64	
}


func NewImage(id string, fileName string, contentFile string, owner string, size float64) *Image {
    return &Image{
        id:          id,
        fileName:    fileName,
        contentFile: contentFile,
        owner:       owner,
        size:        size,
    }
}


func (img *Image) GetId () string{
	return img.id
}

func (img *Image) GetFileName () string{
	return img.fileName
}

func (img *Image) GetContentFile () string{
	return img.contentFile
}

func (img *Image) GetOwner () string{
	return img.owner
}
 
func (img *Image) GetSize () float64{
	return img.size
}

func (img *Image) SetFileName (fileName string){
	img.fileName = fileName	
}

func (img *Image) SetContentFile (contentFile string){
	img.contentFile = contentFile
}

func (img *Image) SetOwner (owner string) {
	img.owner = owner
}
 

