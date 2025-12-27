import { SITE } from '../content/config';

interface Props {
    title?: string;
    description?: string;
    image?: string;
}

export function getSEOMetadata(props: Props, currentPath: URL, site: URL | undefined) {
    const finalTitle = props.title
        ? `${props.title} | ${SITE.title}`
        : SITE.title;

    const finalDesc = props.description || SITE.description;
    const imagePath = props.image || SITE.defaultImage;
    const finalImage = new URL(imagePath, site || "http://localhost:4321"); // sensible default? 
    const canonicalURL = new URL(currentPath.pathname, site || "http://localhost:4321");


    return {
        finalTitle,
        finalDesc,
        finalImage,
        canonicalURL
    };
}